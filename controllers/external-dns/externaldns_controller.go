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
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"kubeops.dev/external-dns-operator/pkg"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ExternalDNSReconciler reconciles a ExternalDNS object
type ExternalDNSReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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

func createAndApplyPlans(edns *externaldnsv1alpha1.ExternalDNS, ctx context.Context) error {
	configs, err := pkg.ConvertCRDtoCfg(*edns)
	if err != nil {
		klog.Info("failed to convert crd into cfg")
		return err
	}

	/*
		// Used for cfg validation
		if err := validation.ValidateConfig(cfg); err != nil {
			klog.Infof("config validation failed: %v", err)
			return
		}
	*/

	for _, cfg := range *configs {
		endpointsSource, err := pkg.CreateEndpointsSource(ctx, &cfg)
		if err != nil {
			klog.Info("failed to create config for domain ", cfg.DomainFilter)
			return err
		}

		provider, err := pkg.CreateProviderFromCfg(ctx, &cfg, endpointsSource)
		if err != nil {
			klog.Info("failed to create provider for domain ", cfg.DomainFilter)
			return err
		}

		reg, err := pkg.CreateRegistry(&cfg, *provider)
		if err != nil {
			klog.Info("failed to create register for domain ", cfg.DomainFilter)
		}

		err = pkg.CreateAndApplySinglePlanForCRD(ctx, &cfg, reg, endpointsSource)
		if err != nil {
			klog.Info("failed to create plan for domain ", cfg.DomainFilter)
			return err
		}
	}
	return nil
}

func (r *ExternalDNSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	//get external dns
	key := req.NamespacedName
	edns := externaldnsv1alpha1.ExternalDNS{}

	//if err := os.Setenv("AWS_DEFAULT_REGION", "us-east-1"); err != nil {
	//	return ctrl.Result{}, err
	//}

	if err := r.Get(ctx, key, &edns); err != nil {
		fmt.Println("failed to get external-dns")
		return ctrl.Result{}, err
	}

	if err := createAndApplyPlans(&edns, ctx); err != nil {
		klog.Info("unable to create plan")
	}

	// dynamic watcher (source service) (later)
	// spec/config function config --> plan.

	// Pending, Current, InProgress

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	///*

	// hasSvc checks whether the external dns crd has a service type source.
	// need to change when the crd changes
	hasSvc := func(ed *externaldnsv1alpha1.ExternalDNS) bool {
		for _, entry := range *ed.Spec.Entries {
			for _, sc := range *entry.Sources.Names {
				if sc == "Service" {
					return true
				}
			}
		}
		return false
	}

	svcHandler := func(object client.Object) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)

		klog.Info("Get service event: ", object.GetName())

		_, found := object.GetAnnotations()["external-dns.alpha.kubernetes.io/hostname"]
		if !found {
			return reconcileReq
		}

		kc := mgr.GetClient()
		dnsList := &externaldnsv1alpha1.ExternalDNSList{}

		if err := kc.List(context.TODO(), dnsList); err != nil {
			klog.Info("failed to list external dns resource")
			return reconcileReq
		}

		for _, ed := range dnsList.Items {
			if hasSvc(&ed) {
				klog.Info("Reconciling service for: ", object.GetName())
				reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: ed.Name, Namespace: ed.Namespace}})
			}
		}
		return reconcileReq
	}

	//*/

	/*
		// for dynamic watcher
		controller, err := ctrl.NewControllerManagedBy(mgr).
			For(&externaldnsv1alpha1.ExternalDNS{}).
			//Watches(pkg.WatchingSources(), handler.EnqueueRequestsFromMapFunc(svcHandler)).
			Build(r)
		if err != nil {
			return err
		}
		// work with the controller
		//controller.Watch(pkg.WatchingSources(), handler.EnqueueRequestsFromMapFunc(svc))
	*/

	return ctrl.NewControllerManagedBy(mgr).
		For(&externaldnsv1alpha1.ExternalDNS{}).
		Watches(pkg.WatchingSources(), handler.EnqueueRequestsFromMapFunc(svcHandler)).
		Complete(r)

}
