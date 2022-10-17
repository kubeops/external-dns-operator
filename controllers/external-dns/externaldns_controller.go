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
	"kubeops.dev/external-dns-operator/pkg/constant"
	"kubeops.dev/external-dns-operator/pkg/credentials"
	"kubeops.dev/external-dns-operator/pkg/informers"
	"kubeops.dev/external-dns-operator/pkg/plan"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ExternalDNSReconciler reconciles a ExternalDNS object
type ExternalDNSReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	watcher *informers.ObjectTracker
}

func (r *ExternalDNSReconciler) GetSecret(ctx context.Context, key types.NamespacedName) (*core.Secret, error) {
	secret := &core.Secret{}
	if err := r.Get(ctx, key, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (r *ExternalDNSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	//get external dns
	ednsKey := req.NamespacedName
	edns := &externaldnsv1alpha1.ExternalDNS{}

	if err := r.Get(ctx, ednsKey, edns); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	edns = edns.DeepCopy()

	edns.Status.Phase = constant.ExternalDNSPhaseInProgress
	_, _, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Phase = constant.ExternalDNSPhaseInProgress
		return in
	})
	if patchErr != nil {
		klog.Error("failed to patch status")
		return ctrl.Result{}, nil
	}

	// dynamic watcher
	if err := informers.RegisterWatcher(ctx, edns, r.watcher, r.Client); err != nil {
		klog.Error("failed to register watcher.", err.Error())
		_, _, patchErr = kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = constant.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(constant.CreateAndRegisterWatcher, "failed to register watcher", in.Generation, false))
			return in
		})
		if patchErr != nil {
			klog.Error("failed to patch status")
			return ctrl.Result{}, patchErr
		}
		return ctrl.Result{}, err
	}

	_, _, patchEr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(constant.CreateAndRegisterWatcher, "watcher registered", in.Generation, true))
		return in
	})
	if patchEr != nil {
		klog.Error("failed to patch status")
		return ctrl.Result{}, patchErr
	}

	//
	// create and set provider secret credentials and environment variables
	secret, err := r.GetSecret(ctx, types.NamespacedName{
		Namespace: edns.Namespace,
		Name:      edns.Spec.ProviderSecretRef.Name})
	if err != nil {
		klog.Error("failed to get provider secret. ", err.Error())

		_, _, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = constant.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(constant.CreateAndSetCredential, err.Error(), in.Generation, false))
			return in
		})
		if patchErr != nil {
			klog.Error("failed to patch status")
		}

		return ctrl.Result{}, patchErr
	}

	err = credentials.SetCredential(secret, ednsKey, edns.Spec.Provider.String())
	if err != nil {

		klog.Error("failed to create credentials. ")

		_, _, patchErr := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = constant.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(constant.CreateAndSetCredential, err.Error(), in.Generation, false))
			return in
		})
		if patchErr != nil {
			klog.Error("failed to patch status")
		}

		return ctrl.Result{}, patchErr
	}

	_, _, patchErr = kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(constant.CreateAndSetCredential, "provider credential configured", in.Generation, true))
		return in
	})
	if patchErr != nil {
		klog.Error("failed to patch status")
		return ctrl.Result{}, patchErr
	}

	//SetDNSRecords creates the dns record according to user information
	//successMsg is used to identify whether the 'plan applied' or 'already up to date'
	successMsg, err := plan.SetDNSRecords(edns, ctx)

	if err != nil {
		klog.Error("failed to create plan. ", err.Error())

		_, _, patchErr = kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = constant.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(constant.CreateAndApplyPlan, err.Error(), in.Generation, false))
			return in
		})
		if patchErr != nil {
			klog.Error("failed to create plan")
		}

		return ctrl.Result{}, patchErr
	}

	if _, _, e := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Phase = constant.ExternalDNSPhaseCurrent
		in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(constant.CreateAndApplyPlan, successMsg, in.Generation, true))
		in.Status.ObservedGeneration = in.Generation
		return in
	}); e != nil {
		klog.Error("failed to patch status")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// for dynamic watcher
	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&externaldnsv1alpha1.ExternalDNS{}).
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
