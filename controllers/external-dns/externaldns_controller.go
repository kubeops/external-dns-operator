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
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	kmc "kmodules.xyz/client-go/client"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"kubeops.dev/external-dns-operator/pkg/credentials"
	"kubeops.dev/external-dns-operator/pkg/informers"
	"kubeops.dev/external-dns-operator/pkg/plan"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ExternalDNSReconciler reconciles a ExternalDNS object
type ExternalDNSReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	watcher *informers.ObjectTracker
}

func (r *ExternalDNSReconciler) getSecret(ctx context.Context, key types.NamespacedName) (*core.Secret, error) {
	secret := &core.Secret{}
	if err := r.Get(ctx, key, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

//update the status of the crd, conditionType is the reason of the condition
func (r *ExternalDNSReconciler) updateEdnsStatus(ctx context.Context, edns *externaldnsv1alpha1.ExternalDNS, conditionType string, conditionMessage string, phase externaldnsv1alpha1.ExternalDNSPhase, conditionStatus bool) error {
	_, _, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Phase = phase
		in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(conditionType, conditionMessage, in.Generation, conditionStatus))
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

	if edns.Status.Phase != externaldnsv1alpha1.ExternalDNSPhaseFailed {
		if patchErr := r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.CreateAndRegisterWatcher, "", externaldnsv1alpha1.ExternalDNSPhaseInProgress, false); patchErr != nil {
			return ctrl.Result{}, patchErr
		}
	}

	// REGISTER WATCHER
	if err := informers.RegisterWatcher(ctx, edns, r.watcher, r.Client); err != nil {
		return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.CreateAndRegisterWatcher, err.Error(), externaldnsv1alpha1.ExternalDNSPhaseFailed, false)
	}

	if patchErr := r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.CreateAndRegisterWatcher, "Watcher Registered", edns.Status.Phase, true); patchErr != nil {
		return ctrl.Result{}, patchErr
	}

	// SECRET AND CREDENTIALS
	// create and set provider secret credentials and environment variables
	secret, err := r.getSecret(ctx, types.NamespacedName{
		Namespace: edns.Namespace,
		Name:      edns.Spec.ProviderSecretRef.Name})
	if err != nil {
		return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.GetProviderSecret, err.Error(), externaldnsv1alpha1.ExternalDNSPhaseFailed, false)
	}

	err = credentials.SetCredential(secret, ednsKey, edns.Spec.Provider.String())
	if err != nil {
		return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.GetProviderSecret, err.Error(), externaldnsv1alpha1.ExternalDNSPhaseFailed, false)
	}

	if patchErr := r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.GetProviderSecret, "provider credential configured", edns.Status.Phase, true); patchErr != nil {
		return ctrl.Result{}, patchErr
	}

	// DNS RECORD
	//SetDNSRecords creates the dns record according to user information
	//successMsg is used to identify whether the 'plan applied' or 'already up to date'
	successMsg, err := plan.SetDNSRecords(edns, ctx)

	if err != nil {
		return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.CreateAndApplyPlan, err.Error(), externaldnsv1alpha1.ExternalDNSPhaseFailed, false)
	}

	return ctrl.Result{}, r.updateEdnsStatus(ctx, edns, externaldnsv1alpha1.CreateAndApplyPlan, successMsg, externaldnsv1alpha1.ExternalDNSPhaseCurrent, true)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	secretToEdns := handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)
		ctx := context.TODO()
		ednsList := &externaldnsv1alpha1.ExternalDNSList{}

		if err := mgr.GetClient().List(ctx, ednsList); err != nil {
			return reconcileReq
		}

		for _, edns := range ednsList.Items {
			if edns.Namespace != object.GetNamespace() || edns.Spec.ProviderSecretRef.Name != object.GetName() {
				continue
			}
			reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
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
