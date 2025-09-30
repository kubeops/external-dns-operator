/*
Copyright AppsCode Inc. and Contributors.

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

package informers

import (
	"context"
	"fmt"
	"reflect"

	api "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	kindNode    = "Node"
	kindService = "Service"
	kindIngress = "Ingress"
)

func getKindNode(cache cache.Cache, r client.Client) (source.SyncingSource, error) {
	hdlr := handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, a *corev1.Node) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)
		dnsList := &api.ExternalDNSList{}

		if err := r.List(ctx, dnsList); err != nil {
			klog.Errorf("failed to list the external dns resources: %s", err.Error())
			return reconcileReq
		}

		for _, edns := range dnsList.Items {
			if edns.Spec.Source.Type.Kind == kindNode {
				reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
			}
		}

		return reconcileReq
	})
	return source.Kind(cache, &corev1.Node{}, hdlr, predicate.TypedFuncs[*corev1.Node]{UpdateFunc: func(e event.TypedUpdateEvent[*corev1.Node]) bool {
		if e.ObjectOld.GetObjectKind().GroupVersionKind().Kind != kindNode {
			return true
		}

		oldNode := e.ObjectOld.DeepCopy()
		newNode := e.ObjectNew.DeepCopy()

		if oldNode.Status.Addresses == nil {
			klog.Error("node addresses not found")
			return false
		}

		if newNode.Status.Addresses == nil {
			klog.Error("node addresses not found")
			return false
		}

		return !reflect.DeepEqual(oldNode.Status.Addresses, newNode.Status.Addresses)
	}}), nil
}

func getKindService(cache cache.Cache, r client.Client) (source.SyncingSource, error) {
	hdlr := handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, a *corev1.Service) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)
		dnsList := &api.ExternalDNSList{}

		if err := r.List(ctx, dnsList); err != nil {
			klog.Errorf("failed to list the external dns resources: %s", err.Error())
			return reconcileReq
		}

		for _, edns := range dnsList.Items {
			if edns.Spec.Source.Type.Kind == kindService {
				reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
			}
		}

		return reconcileReq
	})
	return source.Kind(cache, &corev1.Service{}, hdlr, predicate.TypedFuncs[*corev1.Service]{}), nil
}

func getKindIngress(cache cache.Cache, r client.Client) (source.SyncingSource, error) {
	hdlr := handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, a *networkingv1.Ingress) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)
		dnsList := &api.ExternalDNSList{}

		if err := r.List(ctx, dnsList); err != nil {
			klog.Errorf("failed to list the external dns resources: %s", err.Error())
			return reconcileReq
		}

		for _, edns := range dnsList.Items {
			if edns.Spec.Source.Type.Kind == kindIngress {
				reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
			}
		}

		return reconcileReq
	})

	return source.Kind(cache, &networkingv1.Ingress{}, hdlr, predicate.TypedFuncs[*networkingv1.Ingress]{}), nil
}

func getKind(r client.Client, gvk schema.GroupVersionKind, cache cache.Cache) (source.SyncingSource, error) {
	switch gvk.Kind {
	case kindNode:
		return getKindNode(cache, r)
	case kindService:
		return getKindService(cache, r)
	case kindIngress:
		return getKindIngress(cache, r)
	}
	return nil, fmt.Errorf("unknown kind %v", gvk.Kind)
}
