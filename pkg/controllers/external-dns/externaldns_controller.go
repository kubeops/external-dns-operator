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

	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external/v1alpha1"
	"kubeops.dev/external-dns-operator/pkg/credentials"
	"kubeops.dev/external-dns-operator/pkg/informers"
	"kubeops.dev/external-dns-operator/pkg/plan"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	kmc "kmodules.xyz/client-go/client"
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

func newConditionPtr(reason string, message string, generation int64, conditionStatus bool) *kmapi.Condition {
	newCondition := kmapi.NewCondition(reason, message, generation, conditionStatus)
	return &newCondition
}

func phasePointer(phase externaldnsv1alpha1.ExternalDNSPhase) *externaldnsv1alpha1.ExternalDNSPhase {
	return &phase
}

// update the status of the crd, conditionType is the reason of the condition
func (r *ExternalDNSReconciler) updateEdnsStatus(ctx context.Context, edns *externaldnsv1alpha1.ExternalDNS, newCondition *kmapi.Condition, phase *externaldnsv1alpha1.ExternalDNSPhase) error {
	_, _, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		if phase != nil {
			in.Status.Phase = *phase
		}
		if newCondition != nil {
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, *newCondition)
		}
		return in
	})
	return patchErr
}

func (r ExternalDNSReconciler) patchDNSRecords(ctx context.Context, edns *externaldnsv1alpha1.ExternalDNS, dnsRecs []externaldnsv1alpha1.DNSRecord) error {
	_, _, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.DNSRecords = dnsRecs
		return in
	})

	return patchErr
}

func (r *ExternalDNSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// GET EXTERNAL DNS
	ednsKey := req.NamespacedName
	edns := &externaldnsv1alpha1.ExternalDNS{}

	if err := r.Get(ctx, ednsKey, edns); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	edns = edns.DeepCopy()

	if edns.Status.Phase == "" {
		if patchErr := r.updateEdnsStatus(ctx, edns, nil, phasePointer(externaldnsv1alpha1.ExternalDNSPhaseInProgress)); patchErr != nil {
			return ctrl.Result{}, patchErr
		}
	}

	// REGISTER WATCHER
	if err := informers.RegisterWatcher(ctx, edns, r.watcher, r.Client); err != nil {
		return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, newConditionPtr(externaldnsv1alpha1.CreateAndRegisterWatcher, err.Error(), edns.Generation, false), phasePointer(externaldnsv1alpha1.ExternalDNSPhaseFailed))
	}

	if patchErr := r.updateEdnsStatus(ctx, edns, newConditionPtr(externaldnsv1alpha1.CreateAndRegisterWatcher, "Watcher registered", edns.Generation, true), nil); patchErr != nil {
		return ctrl.Result{}, patchErr
	}

	mutex.Lock()
	defer mutex.Unlock()

	// SECRET AND CREDENTIALS
	// create and set provider secret credentials and environment variables
	err := credentials.SetCredential(ctx, r.Client, edns)
	if err != nil {
		return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, newConditionPtr(externaldnsv1alpha1.GetProviderSecret, err.Error(), edns.Generation, false), phasePointer(externaldnsv1alpha1.ExternalDNSPhaseFailed))
	}

	if patchErr := r.updateEdnsStatus(ctx, edns, newConditionPtr(externaldnsv1alpha1.GetProviderSecret, "Provider credential configured", edns.Generation, true), nil); patchErr != nil {
		return ctrl.Result{}, patchErr
	}

	// APPLY DNS RECORD
	// SetDNSRecords creates the dns record according to user information
	// successMsg is used to identify whether the 'plan applied' or 'already up to date'
	dnsRecs, err := plan.SetDNSRecords(ctx, edns)
	if err != nil {
		return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, newConditionPtr(externaldnsv1alpha1.CreateAndApplyPlan, err.Error(), edns.Generation, false), phasePointer(externaldnsv1alpha1.ExternalDNSPhaseFailed))
	}

	err = r.patchDNSRecords(ctx, edns, dnsRecs)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, newConditionPtr(externaldnsv1alpha1.CreateAndApplyPlan, "plan applied", edns.Generation, true), phasePointer(externaldnsv1alpha1.ExternalDNSPhaseCurrent))
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	secretToEdns := handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)
		ctx := context.TODO()
		ednsList := &externaldnsv1alpha1.ExternalDNSList{}

		if err := mgr.GetClient().List(ctx, ednsList, client.InNamespace(object.GetNamespace())); err != nil {
			return reconcileReq
		}

		for _, edns := range ednsList.Items {
			if edns.Spec.ProviderSecretRef != nil && edns.Spec.ProviderSecretRef.Name == object.GetName() {
				reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
			}
		}

		return reconcileReq
	})

	// for dynamic watcher
	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&externaldnsv1alpha1.ExternalDNS{}).
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
