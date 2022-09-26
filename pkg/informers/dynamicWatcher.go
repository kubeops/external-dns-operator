package informers

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync"
)

type ObjectTracker struct {
	m sync.Map

	controller.Controller
}

func (o ObjectTracker) Watch(obj runtime.Object, handler handler.EventHandler) error {
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
	)
	if err != nil {
		o.m.Delete(key)
		return errors.Wrapf(err, "failed to add watcher on external object %q", gvk.String())
	}
	return nil
}

func sourceHandler(object client.Object) []reconcile.Request {

}

func getRuntimeObject(gvk schema.GroupVersionKind) (runtime.Object, error) {
	unObj := &unstructured.Unstructured{}
	unObj.SetGroupVersionKind(gvk)
	return unObj, nil
	//if *src == "service"{
	//	return &v1.Service{},nil
	//}
	//if *src == "ingress" {
	//
	//}

}

func RegisterWatcher(sourceList []string, watcher ObjectTracker) {
	for _, source := range sourceList {

		watcher.Watch(&unstructured.Unstructured{}, handler.EnqueueRequestsFromMapFunc(sourceHandler))
	}
}
