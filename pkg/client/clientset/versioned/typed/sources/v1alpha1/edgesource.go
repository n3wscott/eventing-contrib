/*
Copyright 2018 The Knative Authors

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
	v1alpha1 "github.com/knative/eventing-sources/pkg/apis/sources/v1alpha1"
	scheme "github.com/knative/eventing-sources/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EdgeSourcesGetter has a method to return a EdgeSourceInterface.
// A group's client should implement this interface.
type EdgeSourcesGetter interface {
	EdgeSources(namespace string) EdgeSourceInterface
}

// EdgeSourceInterface has methods to work with EdgeSource resources.
type EdgeSourceInterface interface {
	Create(*v1alpha1.EdgeSource) (*v1alpha1.EdgeSource, error)
	Update(*v1alpha1.EdgeSource) (*v1alpha1.EdgeSource, error)
	UpdateStatus(*v1alpha1.EdgeSource) (*v1alpha1.EdgeSource, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.EdgeSource, error)
	List(opts v1.ListOptions) (*v1alpha1.EdgeSourceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EdgeSource, err error)
	EdgeSourceExpansion
}

// edgeSources implements EdgeSourceInterface
type edgeSources struct {
	client rest.Interface
	ns     string
}

// newEdgeSources returns a EdgeSources
func newEdgeSources(c *SourcesV1alpha1Client, namespace string) *edgeSources {
	return &edgeSources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the edgeSource, and returns the corresponding edgeSource object, and an error if there is any.
func (c *edgeSources) Get(name string, options v1.GetOptions) (result *v1alpha1.EdgeSource, err error) {
	result = &v1alpha1.EdgeSource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("edgesources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EdgeSources that match those selectors.
func (c *edgeSources) List(opts v1.ListOptions) (result *v1alpha1.EdgeSourceList, err error) {
	result = &v1alpha1.EdgeSourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("edgesources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested edgeSources.
func (c *edgeSources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("edgesources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a edgeSource and creates it.  Returns the server's representation of the edgeSource, and an error, if there is any.
func (c *edgeSources) Create(edgeSource *v1alpha1.EdgeSource) (result *v1alpha1.EdgeSource, err error) {
	result = &v1alpha1.EdgeSource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("edgesources").
		Body(edgeSource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a edgeSource and updates it. Returns the server's representation of the edgeSource, and an error, if there is any.
func (c *edgeSources) Update(edgeSource *v1alpha1.EdgeSource) (result *v1alpha1.EdgeSource, err error) {
	result = &v1alpha1.EdgeSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("edgesources").
		Name(edgeSource.Name).
		Body(edgeSource).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *edgeSources) UpdateStatus(edgeSource *v1alpha1.EdgeSource) (result *v1alpha1.EdgeSource, err error) {
	result = &v1alpha1.EdgeSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("edgesources").
		Name(edgeSource.Name).
		SubResource("status").
		Body(edgeSource).
		Do().
		Into(result)
	return
}

// Delete takes name of the edgeSource and deletes it. Returns an error if one occurs.
func (c *edgeSources) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("edgesources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *edgeSources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("edgesources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched edgeSource.
func (c *edgeSources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EdgeSource, err error) {
	result = &v1alpha1.EdgeSource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("edgesources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
