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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "kubeops.dev/external-dns-operator/apis/external/v1alpha1"
	scheme "kubeops.dev/external-dns-operator/client/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ExternalDNSsGetter has a method to return a ExternalDNSInterface.
// A group's client should implement this interface.
type ExternalDNSsGetter interface {
	ExternalDNSs(namespace string) ExternalDNSInterface
}

// ExternalDNSInterface has methods to work with ExternalDNS resources.
type ExternalDNSInterface interface {
	Create(ctx context.Context, externalDNS *v1alpha1.ExternalDNS, opts v1.CreateOptions) (*v1alpha1.ExternalDNS, error)
	Update(ctx context.Context, externalDNS *v1alpha1.ExternalDNS, opts v1.UpdateOptions) (*v1alpha1.ExternalDNS, error)
	UpdateStatus(ctx context.Context, externalDNS *v1alpha1.ExternalDNS, opts v1.UpdateOptions) (*v1alpha1.ExternalDNS, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ExternalDNS, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ExternalDNSList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ExternalDNS, err error)
	ExternalDNSExpansion
}

// externalDNSs implements ExternalDNSInterface
type externalDNSs struct {
	client rest.Interface
	ns     string
}

// newExternalDNSs returns a ExternalDNSs
func newExternalDNSs(c *ExternalDNSV1alpha1Client, namespace string) *externalDNSs {
	return &externalDNSs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the externalDNS, and returns the corresponding externalDNS object, and an error if there is any.
func (c *externalDNSs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ExternalDNS, err error) {
	result = &v1alpha1.ExternalDNS{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("externaldnss").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ExternalDNSs that match those selectors.
func (c *externalDNSs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ExternalDNSList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ExternalDNSList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("externaldnss").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested externalDNSs.
func (c *externalDNSs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("externaldnss").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a externalDNS and creates it.  Returns the server's representation of the externalDNS, and an error, if there is any.
func (c *externalDNSs) Create(ctx context.Context, externalDNS *v1alpha1.ExternalDNS, opts v1.CreateOptions) (result *v1alpha1.ExternalDNS, err error) {
	result = &v1alpha1.ExternalDNS{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("externaldnss").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(externalDNS).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a externalDNS and updates it. Returns the server's representation of the externalDNS, and an error, if there is any.
func (c *externalDNSs) Update(ctx context.Context, externalDNS *v1alpha1.ExternalDNS, opts v1.UpdateOptions) (result *v1alpha1.ExternalDNS, err error) {
	result = &v1alpha1.ExternalDNS{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("externaldnss").
		Name(externalDNS.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(externalDNS).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *externalDNSs) UpdateStatus(ctx context.Context, externalDNS *v1alpha1.ExternalDNS, opts v1.UpdateOptions) (result *v1alpha1.ExternalDNS, err error) {
	result = &v1alpha1.ExternalDNS{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("externaldnss").
		Name(externalDNS.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(externalDNS).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the externalDNS and deletes it. Returns an error if one occurs.
func (c *externalDNSs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("externaldnss").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *externalDNSs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("externaldnss").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched externalDNS.
func (c *externalDNSs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ExternalDNS, err error) {
	result = &v1alpha1.ExternalDNS{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("externaldnss").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
