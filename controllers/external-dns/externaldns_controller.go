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
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"kubeops.dev/external-dns-operator/pkg"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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

func createPlans(edns externaldnsv1alpha1.ExternalDNS, ctx *context.Context) {
	cfg, err := pkg.ConvertCRDtoCfg(edns)
	if err != nil {
		fmt.Println("failed to convert crd into cfg")
		return
	}

	endpointsSource, err := pkg.CreateEndpointsSource(cfg, ctx)
	if err != nil {
		fmt.Println("failed to create endpoints source")
		return
	}

	provider, err := pkg.CreateProviderFromCfg(cfg, *ctx, endpointsSource)
	if err != nil {
		fmt.Println("failed to create provider")
	}

	reg, err := pkg.CreateRegistry(cfg, *provider)
	if err != nil {
		fmt.Println("failed to create registry")
		return
	}

	plan, err := pkg.CreateSinglePlanForCRD(cfg, reg, *ctx, *endpointsSource)
	if err != nil {
		fmt.Println("failed to create plan")
		return
	}

	plan = plan.Calculate()

	if plan.Changes.HasChanges() {
		err = reg.ApplyChanges(*ctx, plan.Changes)
		if err != nil {
			fmt.Println("failed to apply changes for plan")
			return
		}
	} else {
		fmt.Println("all records are already up to date")
	}

	return
}

func (r *ExternalDNSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	// get external dns

	key := req.NamespacedName
	edns := externaldnsv1alpha1.ExternalDNS{}

	if err := r.Get(ctx, key, &edns); err != nil {
		fmt.Println("failed to get external-dns")
		return ctrl.Result{}, err
	}
	fmt.Println("found external-dns")
	fmt.Println("Zone : ", edns.Spec.AWSZone)

	// pkg/provider/aws.go
	// dynamic watcher (source service) (later)
	// spec/config function config --> plan.

	// Pending, Current, InProgress

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&externaldnsv1alpha1.ExternalDNS{}).
		//Watches() // service, handler
		Complete(r)
}
