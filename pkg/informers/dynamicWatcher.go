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
	"sync"

	api "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type ObjectTracker struct {
	m sync.Map

	Manager manager.Manager
	controller.Controller
}

func (o *ObjectTracker) Watch(obj runtime.Object, r client.Client) error {
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

	kind, err := getKind(r, gvk, o.Manager.GetCache())
	if err != nil {
		klog.Error(err, "unable to watch object "+gvk.String())
		return err
	}
	if err = o.Controller.Watch(kind); err != nil {
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
	return watcher.Watch(getRuntimeObject(crd.Spec.Source.Type.GroupVersionKind()), r)
}
