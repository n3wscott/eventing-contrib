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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	sources_v1alpha1 "github.com/knative/eventing-sources/contrib/gcppubsub/pkg/apis/sources/v1alpha1"
	versioned "github.com/knative/eventing-sources/contrib/gcppubsub/pkg/client/clientset/versioned"
	internalinterfaces "github.com/knative/eventing-sources/contrib/gcppubsub/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/knative/eventing-sources/contrib/gcppubsub/pkg/client/listers/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GcpPubSubSourceInformer provides access to a shared informer and lister for
// GcpPubSubSources.
type GcpPubSubSourceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.GcpPubSubSourceLister
}

type gcpPubSubSourceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewGcpPubSubSourceInformer constructs a new informer for GcpPubSubSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGcpPubSubSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGcpPubSubSourceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredGcpPubSubSourceInformer constructs a new informer for GcpPubSubSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGcpPubSubSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SourcesV1alpha1().GcpPubSubSources(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SourcesV1alpha1().GcpPubSubSources(namespace).Watch(options)
			},
		},
		&sources_v1alpha1.GcpPubSubSource{},
		resyncPeriod,
		indexers,
	)
}

func (f *gcpPubSubSourceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGcpPubSubSourceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gcpPubSubSourceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&sources_v1alpha1.GcpPubSubSource{}, f.defaultInformer)
}

func (f *gcpPubSubSourceInformer) Lister() v1alpha1.GcpPubSubSourceLister {
	return v1alpha1.NewGcpPubSubSourceLister(f.Informer().GetIndexer())
}
