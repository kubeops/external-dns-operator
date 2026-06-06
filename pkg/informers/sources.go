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
		// Reconcile only when the addresses, labels, or annotations change —
		// these are the only fields that affect the DNS endpoints external-dns
		// derives from a Node.
		return !reflect.DeepEqual(e.ObjectOld.Status.Addresses, e.ObjectNew.Status.Addresses) ||
			!reflect.DeepEqual(e.ObjectOld.Labels, e.ObjectNew.Labels) ||
			!reflect.DeepEqual(e.ObjectOld.Annotations, e.ObjectNew.Annotations)
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
	return source.Kind(cache, &corev1.Service{}, hdlr, predicate.TypedFuncs[*corev1.Service]{UpdateFunc: func(e event.TypedUpdateEvent[*corev1.Service]) bool {
		// Reconcile only on changes that affect the endpoints external-dns
		// derives from a Service: spec changes (bumps Generation), annotations
		// (hostname/ttl/etc.), labels (selector matching), or the LB ingress
		// status. This skips no-op resync events.
		return e.ObjectOld.Generation != e.ObjectNew.Generation ||
			!reflect.DeepEqual(e.ObjectOld.Annotations, e.ObjectNew.Annotations) ||
			!reflect.DeepEqual(e.ObjectOld.Labels, e.ObjectNew.Labels) ||
			!reflect.DeepEqual(e.ObjectOld.Status.LoadBalancer, e.ObjectNew.Status.LoadBalancer)
	}}), nil
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

	return source.Kind(cache, &networkingv1.Ingress{}, hdlr, predicate.TypedFuncs[*networkingv1.Ingress]{UpdateFunc: func(e event.TypedUpdateEvent[*networkingv1.Ingress]) bool {
		// Reconcile only on changes that affect the endpoints external-dns
		// derives from an Ingress: spec changes (bumps Generation),
		// annotations (hostname/ttl/etc.), labels, or the LB ingress status.
		return e.ObjectOld.Generation != e.ObjectNew.Generation ||
			!reflect.DeepEqual(e.ObjectOld.Annotations, e.ObjectNew.Annotations) ||
			!reflect.DeepEqual(e.ObjectOld.Labels, e.ObjectNew.Labels) ||
			!reflect.DeepEqual(e.ObjectOld.Status.LoadBalancer, e.ObjectNew.Status.LoadBalancer)
	}}), nil
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
