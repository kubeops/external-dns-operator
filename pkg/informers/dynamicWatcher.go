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
	"reflect"
	"sync"

	api "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type ObjectTracker struct {
	m sync.Map

	Manager manager.Manager
	controller.Controller
}

func (o *ObjectTracker) Watch(obj runtime.Object, handler handler.EventHandler) error {
	if o.Controller == nil {
		return nil
	}

	gvk := obj.GetObjectKind().GroupVersionKind()
	key := gvk.GroupKind().String()

	if _, loaded := o.m.LoadOrStore(key, struct{}{}); loaded {
		return nil
	}

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(gvk)

	// adding watcher to an external object
	err := o.Controller.Watch(
		source.Kind(o.Manager.GetCache(), u),
		handler,
		predicate.Funcs{UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld.GetObjectKind().GroupVersionKind().Kind != "Node" {
				return true
			}

			oldNode := e.ObjectOld.(*unstructured.Unstructured).DeepCopy()
			newNode := e.ObjectNew.(*unstructured.Unstructured).DeepCopy()

			oldAddr, found, err := unstructured.NestedSlice(oldNode.Object, "status", "addresses")
			if err != nil {
				klog.Error(err.Error())
				return false
			}

			if !found {
				klog.Error("can't found node addresses")
				return false
			}

			newAddr, found, err := unstructured.NestedSlice(newNode.Object, "status", "addresses")
			if err != nil {
				klog.Error(err.Error())
				return false
			}

			if !found {
				klog.Error("can't found node addresses")
				return false
			}

			return !reflect.DeepEqual(oldAddr, newAddr)
		}},
	)
	if err != nil {
		o.m.Delete(key)
		return errors.Wrapf(err, "failed to add watcher on external object %q", gvk.String())
	}
	return nil
}

func getRuntimeObject(gvk schema.GroupVersionKind) runtime.Object {
	unObj := &unstructured.Unstructured{}
	unObj.SetGroupVersionKind(gvk)
	return unObj
}

func RegisterWatcher(ctx context.Context, crd *api.ExternalDNS, watcher *ObjectTracker, r client.Client) error {
	sourceHandler := func(ctx context.Context, object client.Object) []reconcile.Request {
		reconcileReq := make([]reconcile.Request, 0)

		dnsList := &api.ExternalDNSList{}

		if err := r.List(ctx, dnsList); err != nil {
			klog.Errorf("failed to list the external dns resources: %s", err.Error())
			return reconcileReq
		}

		objKind := object.GetObjectKind().GroupVersionKind().Kind

		for _, edns := range dnsList.Items {
			if edns.Spec.Source.Type.Kind == objKind {
				reconcileReq = append(reconcileReq, reconcile.Request{NamespacedName: client.ObjectKey{Name: edns.Name, Namespace: edns.Namespace}})
			}
		}

		return reconcileReq
	}

	return watcher.Watch(getRuntimeObject(crd.Spec.Source.Type.GroupVersionKind()), handler.EnqueueRequestsFromMapFunc(sourceHandler))
}
