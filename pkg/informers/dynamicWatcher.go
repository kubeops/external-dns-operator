package informers

import (
	"context"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync"
)

type ObjectTracker struct {
	m sync.Map

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
		&source.Kind{Type: u},
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

func RegisterWatcher(ctx context.Context, crd *externaldnsv1alpha1.ExternalDNS, watcher *ObjectTracker, r client.Client) error {

	sourceHandler := func(object client.Object) []reconcile.Request {

		reconcileReq := make([]reconcile.Request, 0)

		dnsList := &externaldnsv1alpha1.ExternalDNSList{}

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
