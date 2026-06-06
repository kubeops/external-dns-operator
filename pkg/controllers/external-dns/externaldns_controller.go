/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package externaldns

import (
	"context"
	"sync"

	api "kubeops.dev/external-dns-operator/apis/external/v1alpha1"
	"kubeops.dev/external-dns-operator/pkg/credentials"
	"kubeops.dev/external-dns-operator/pkg/informers"
	"kubeops.dev/external-dns-operator/pkg/plan"

	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	kmc "kmodules.xyz/client-go/client"
	condutil "kmodules.xyz/client-go/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var mutex sync.Mutex

const finalizer = "externaldns.kubeops.dev/finalizer"

// ExternalDNSReconciler reconciles a ExternalDNS object
type ExternalDNSReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	watcher *informers.ObjectTracker
}

func newCondition(reason string, message string, generation int64, conditionStatus bool) *kmapi.Condition {
	newCondition := condutil.NewCondition(reason, message, generation, conditionStatus)
	return &newCondition
}

func newPhase(phase api.ExternalDNSPhase) *api.ExternalDNSPhase {
	return &phase
}

// statusUpdate accumulates all fields the reconciler may want to change
// on Status. The reconciler builds one of these per Reconcile pass and
// flushes it with a single PatchStatus call instead of issuing 4-5
// independent patches.
type statusUpdate struct {
	conditions []kmapi.Condition
	phase      *api.ExternalDNSPhase
	dnsRecords *[]api.DNSRecord
}

func (s *statusUpdate) setCondition(c *kmapi.Condition) *statusUpdate {
	if c != nil {
		s.conditions = append(s.conditions, *c)
	}
	return s
}

func (s *statusUpdate) setPhase(p *api.ExternalDNSPhase) *statusUpdate {
	s.phase = p
	return s
}

func (s *statusUpdate) setDNSRecords(recs []api.DNSRecord) *statusUpdate {
	s.dnsRecords = &recs
	return s
}

func (r *ExternalDNSReconciler) patchStatus(ctx context.Context, edns *api.ExternalDNS, update *statusUpdate) error {
	generation := edns.Generation
	_, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*api.ExternalDNS)
		in.Status.ObservedGeneration = generation
		if update.phase != nil {
			in.Status.Phase = *update.phase
		}
		for i := range update.conditions {
			in.Status.Conditions = condutil.SetCondition(in.Status.Conditions, update.conditions[i])
		}
		if update.dnsRecords != nil {
			in.Status.DNSRecords = *update.dnsRecords
		}
		return in
	})
	return patchErr
}

func (r *ExternalDNSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// GET EXTERNAL DNS
	ednsKey := req.NamespacedName
	edns := &api.ExternalDNS{}

	if err := r.Get(ctx, ednsKey, edns); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	edns = edns.DeepCopy()

	// HANDLE DELETION
	if !edns.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, edns)
	}

	// ENSURE FINALIZER
	if !controllerutil.ContainsFinalizer(edns, finalizer) {
		controllerutil.AddFinalizer(edns, finalizer)
		if err := r.Update(ctx, edns); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	update := &statusUpdate{}
	if edns.Status.Phase == "" {
		update.setPhase(newPhase(api.ExternalDNSPhaseInProgress))
	}

	// REGISTER WATCHER
	if err := informers.RegisterWatcher(ctx, edns, r.watcher, r.Client); err != nil {
		update.
			setCondition(newCondition(api.CreateAndRegisterWatcher, err.Error(), edns.Generation, false)).
			setPhase(newPhase(api.ExternalDNSPhaseFailed))
		if patchErr := r.patchStatus(ctx, edns, update); patchErr != nil {
			err = errors.Wrap(err, patchErr.Error())
		}
		return ctrl.Result{}, err
	}
	update.setCondition(newCondition(api.CreateAndRegisterWatcher, "Watcher registered", edns.Generation, true))

	mutex.Lock()
	defer mutex.Unlock()

	// SECRET AND CREDENTIALS
	if err := credentials.SetCredential(ctx, r.Client, edns); err != nil {
		update.
			setCondition(newCondition(api.GetProviderSecret, err.Error(), edns.Generation, false)).
			setPhase(newPhase(api.ExternalDNSPhaseFailed))
		if patchErr := r.patchStatus(ctx, edns, update); patchErr != nil {
			err = errors.Wrap(err, patchErr.Error())
		}
		return ctrl.Result{}, err
	}
	update.setCondition(newCondition(api.GetProviderSecret, "Provider credential configured", edns.Generation, true))

	// APPLY DNS RECORD
	dnsRecs, err := plan.SetDNSRecords(ctx, edns)
	if err != nil {
		update.
			setCondition(newCondition(api.CreateAndApplyPlan, err.Error(), edns.Generation, false)).
			setPhase(newPhase(api.ExternalDNSPhaseFailed))
		if patchErr := r.patchStatus(ctx, edns, update); patchErr != nil {
			err = errors.Wrap(err, patchErr.Error())
		}
		return ctrl.Result{}, err
	}

	update.setDNSRecords(dnsRecs)
	if len(dnsRecs) == 0 {
		update.
			setCondition(newCondition(api.CreateAndApplyPlan, "no endpoints found for source", edns.Generation, true)).
			setPhase(newPhase(api.ExternalDNSPhaseInProgress))
	} else {
		update.
			setCondition(newCondition(api.CreateAndApplyPlan, "plan applied", edns.Generation, true)).
			setPhase(newPhase(api.ExternalDNSPhaseCurrent))
	}
	return ctrl.Result{}, r.patchStatus(ctx, edns, update)
}

func (r *ExternalDNSReconciler) handleDeletion(ctx context.Context, edns *api.ExternalDNS) (ctrl.Result, error) {
	if !controllerutil.ContainsFinalizer(edns, finalizer) {
		return ctrl.Result{}, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	if edns.Spec.Policy == nil || *edns.Spec.Policy == api.PolicySync {
		if err := credentials.SetCredential(ctx, r.Client, edns); err != nil {
			klog.Errorf("failed to set credentials for cleanup of %s/%s: %v", edns.Namespace, edns.Name, err)
			return ctrl.Result{}, err
		}

		if err := plan.DeleteDNSRecords(ctx, edns); err != nil {
			klog.Errorf("failed to delete DNS records for %s/%s: %v", edns.Namespace, edns.Name, err)
			return ctrl.Result{}, err
		}
	}

	controllerutil.RemoveFinalizer(edns, finalizer)
	return ctrl.Result{}, r.Update(ctx, edns)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	secretToEdns := handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)
		ednsList := &api.ExternalDNSList{}

		if err := mgr.GetClient().List(ctx, ednsList, client.InNamespace(object.GetNamespace())); err != nil {
			return reconcileReq
		}

		for _, edns := range ednsList.Items {
			switch edns.Spec.Provider {
			case api.ProviderAWS:
				if edns.Spec.AWS != nil && edns.Spec.AWS.SecretRef != nil && edns.Spec.AWS.SecretRef.Name == object.GetName() {
					reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
				}

			case api.ProviderAzure:
				if edns.Spec.Azure != nil && edns.Spec.Azure.SecretRef != nil && edns.Spec.Azure.SecretRef.Name == object.GetName() {
					reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
				}

			case api.ProviderGoogle:
				if edns.Spec.Google != nil && edns.Spec.Google.SecretRef != nil && edns.Spec.Google.SecretRef.Name == object.GetName() {
					reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
				}

			case api.ProviderCloudflare:
				if edns.Spec.Cloudflare != nil && edns.Spec.Cloudflare.SecretRef != nil && edns.Spec.Cloudflare.SecretRef.Name == object.GetName() {
					reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
				}
			}
		}

		return reconcileReq
	})

	// for dynamic watcher
	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&api.ExternalDNS{}).
		Watches(&core.Secret{}, secretToEdns).
		Build(r)
	if err != nil {
		klog.Error("failed to build controller.", err.Error())
		return err
	}

	r.watcher = &informers.ObjectTracker{
		Manager:    mgr,
		Controller: controller,
	}

	return nil
}
