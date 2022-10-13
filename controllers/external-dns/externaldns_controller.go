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
	v1 "k8s.io/api/core/v1"
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
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ExternalDNSReconciler reconciles a ExternalDNS object
type ExternalDNSReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	watcher *informers.ObjectTracker
}

//+kubebuilder:rbac:groups=external-dns.appscode.com,resources=externaldns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=external-dns.appscode.com,resources=externaldns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=external-dns.appscode.com,resources=externaldns/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ExternalDNS object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile

func (r *ExternalDNSReconciler) GetSecret(ctx context.Context, key *types.NamespacedName) (*v1.Secret, error) {
	secret := &v1.Secret{}
	if err := r.Get(ctx, *key, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (r *ExternalDNSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	//get external dns
	key := req.NamespacedName
	edns := &externaldnsv1alpha1.ExternalDNS{}

	if err := r.Get(ctx, key, edns); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	edns = edns.DeepCopy()

	edns.Status.Phase = externaldnsv1alpha1.ExternalDNSPhaseInProgress

	// dynamic watcher
	if err := informers.RegisterWatcher(ctx, edns, r.watcher, r.Client); err != nil {
		klog.Info("failed to register watcher.", err.Error())
		_, _, err = kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = externaldnsv1alpha1.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(externaldnsv1alpha1.ConditionWatcher, "failed to register watcher", in.Generation, false))
			return in
		})
		klog.Info("failed to patch status")
		return ctrl.Result{}, err
	}
	if _, _, err := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(externaldnsv1alpha1.ConditionWatcher, "watcher registered", in.Generation, true))
		return in
	}); err != nil {
		klog.Info("failed to patch status")
	}

	//
	// Provider secret credentials and environment variables
	secret, err := r.GetSecret(ctx, &types.NamespacedName{
		Namespace: edns.Namespace,
		Name:      edns.Spec.ProviderSecretRef.Name})
	if err != nil {
		klog.Info("failed to get provider secret. ", err.Error())

		if _, _, e := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = externaldnsv1alpha1.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(externaldnsv1alpha1.ConditionCredential, "failed to create/set provider credential", in.Generation, false))
			return in
		}); e != nil {
			klog.Info("failed to patch status")
		}

		return ctrl.Result{}, err
	}

	err = credentials.SetCredential(secret, key, edns.Spec.Provider.String())
	if err != nil {

		if _, _, e := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = externaldnsv1alpha1.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(externaldnsv1alpha1.ConditionCredential, "failed to create/set provider credential", in.Generation, false))
			return in
		}); e != nil {
			klog.Info("failed to patch status")
		}

		return ctrl.Result{}, err
	}

	if _, _, e := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(externaldnsv1alpha1.ConditionCredential, "provider credential configured", in.Generation, true))
		return in
	}); e != nil {
		klog.Info("failed to patch status")
	}

	//Related to making plan
	//successMsg is used to identify whether the 'plan applied' or 'already up to date'
	successMsg, err := plan.MakePlan(edns, ctx)

	if err != nil {
		klog.Info("failed to create plan. ", err.Error())

		if _, _, e := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
			in := obj.(*externaldnsv1alpha1.ExternalDNS)
			in.Status.Phase = externaldnsv1alpha1.ExternalDNSPhaseFailed
			in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(externaldnsv1alpha1.ConditionPlan, "failed to create/apply plan", in.Generation, false))
			return in
		}); e != nil {
			klog.Info("failed to create plan")
		}

		return ctrl.Result{}, err
	}

	if _, _, e := kmc.PatchStatus(ctx, r.Client, edns, func(obj client.Object) client.Object {
		in := obj.(*externaldnsv1alpha1.ExternalDNS)
		in.Status.Phase = externaldnsv1alpha1.ExternalDNSPhaseCurrent
		in.Status.Conditions = kmapi.SetCondition(in.Status.Conditions, kmapi.NewCondition(externaldnsv1alpha1.ConditionPlan, *successMsg, in.Generation, true))
		in.Status.ObservedGeneration = in.Generation
		return in
	}); e != nil {
		klog.Info("failed to patch status")
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
		klog.Infof("failed to build controller.", err.Error())
		return err
	}

	r.watcher = &informers.ObjectTracker{
		Controller: controller,
	}

	return nil
}
