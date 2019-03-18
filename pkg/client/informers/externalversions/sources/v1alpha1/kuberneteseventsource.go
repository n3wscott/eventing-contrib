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

	sources_v1alpha1 "github.com/knative/eventing-sources/pkg/apis/sources/v1alpha1"
	versioned "github.com/knative/eventing-sources/pkg/client/clientset/versioned"
	internalinterfaces "github.com/knative/eventing-sources/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/knative/eventing-sources/pkg/client/listers/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// KubernetesEventSourceInformer provides access to a shared informer and lister for
// KubernetesEventSources.
type KubernetesEventSourceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.KubernetesEventSourceLister
}

type kubernetesEventSourceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewKubernetesEventSourceInformer constructs a new informer for KubernetesEventSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewKubernetesEventSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredKubernetesEventSourceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredKubernetesEventSourceInformer constructs a new informer for KubernetesEventSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredKubernetesEventSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SourcesV1alpha1().KubernetesEventSources(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SourcesV1alpha1().KubernetesEventSources(namespace).Watch(options)
			},
		},
		&sources_v1alpha1.KubernetesEventSource{},
		resyncPeriod,
		indexers,
	)
}

func (f *kubernetesEventSourceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredKubernetesEventSourceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *kubernetesEventSourceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&sources_v1alpha1.KubernetesEventSource{}, f.defaultInformer)
}

func (f *kubernetesEventSourceInformer) Lister() v1alpha1.KubernetesEventSourceLister {
	return v1alpha1.NewKubernetesEventSourceLister(f.Informer().GetIndexer())
}
