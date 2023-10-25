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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var mutex sync.Mutex

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

// update the status of the crd, conditionType is the reason of the condition
func (r *ExternalDNSReconciler) updateEdnsStatus(ctx context.Context, edns *api.ExternalDNS, newCondition *kmapi.Condition, phase *api.ExternalDNSPhase) error {
	_, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*api.ExternalDNS)
		if phase != nil {
			in.Status.Phase = *phase
		}
		if newCondition != nil {
			in.Status.Conditions = condutil.SetCondition(in.Status.Conditions, *newCondition)
		}
		return in
	})
	return patchErr
}

func (r ExternalDNSReconciler) patchDNSRecords(ctx context.Context, edns *api.ExternalDNS, dnsRecs []api.DNSRecord) error {
	_, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*api.ExternalDNS)
		in.Status.DNSRecords = dnsRecs
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

	if edns.Status.Phase == "" {
		if patchErr := r.updateEdnsStatus(
			ctx,
			edns,
			nil,
			newPhase(api.ExternalDNSPhaseInProgress),
		); patchErr != nil {
			return ctrl.Result{}, patchErr
		}
	}

	// REGISTER WATCHER
	if err := informers.RegisterWatcher(ctx, edns, r.watcher, r.Client); err != nil {
		if patchErr := r.updateEdnsStatus(
			ctx,
			edns,
			newCondition(api.CreateAndRegisterWatcher, err.Error(), edns.Generation, false),
			newPhase(api.ExternalDNSPhaseFailed),
		); patchErr != nil {
			err = errors.Wrap(err, patchErr.Error())
		}
		return ctrl.Result{}, err
	}

	if patchErr := r.updateEdnsStatus(
		ctx,
		edns,
		newCondition(api.CreateAndRegisterWatcher, "Watcher registered", edns.Generation, true),
		nil,
	); patchErr != nil {
		return ctrl.Result{}, patchErr
	}

	mutex.Lock()
	defer mutex.Unlock()

	// SECRET AND CREDENTIALS
	// create and set provider secret credentials and environment variables
	err := credentials.SetCredential(ctx, r.Client, edns)
	if err != nil {
		if patchErr := r.updateEdnsStatus(
			ctx,
			edns,
			newCondition(api.GetProviderSecret, err.Error(), edns.Generation, false),
			newPhase(api.ExternalDNSPhaseFailed),
		); patchErr != nil {
			err = errors.Wrap(err, patchErr.Error())
		}
		return ctrl.Result{}, err
	}

	if patchErr := r.updateEdnsStatus(
		ctx,
		edns,
		newCondition(api.GetProviderSecret, "Provider credential configured", edns.Generation, true),
		nil,
	); patchErr != nil {
		return ctrl.Result{}, patchErr
	}

	// APPLY DNS RECORD
	// SetDNSRecords creates the dns record according to user information
	// successMsg is used to identify whether the 'plan applied' or 'already up to date'
	dnsRecs, err := plan.SetDNSRecords(ctx, edns)
	if err != nil {
		if patchErr := r.updateEdnsStatus(
			ctx,
			edns,
			newCondition(api.CreateAndApplyPlan, err.Error(), edns.Generation, false),
			newPhase(api.ExternalDNSPhaseFailed),
		); patchErr != nil {
			err = errors.Wrap(err, patchErr.Error())
		}
		return ctrl.Result{}, err
	}

	err = r.patchDNSRecords(ctx, edns, dnsRecs)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, r.updateEdnsStatus(
		ctx,
		edns,
		newCondition(api.CreateAndApplyPlan, "plan applied", edns.Generation, true),
		newPhase(api.ExternalDNSPhaseCurrent),
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	secretToEdns := handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)
		ctx := context.TODO()
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
		Watches(&source.Kind{Type: &core.Secret{}}, secretToEdns).
		Build(r)
	if err != nil {
		klog.Error("failed to build controller.", err.Error())
		return err
	}

	r.watcher = &informers.ObjectTracker{
		Controller: controller,
	}

	return nil
}
