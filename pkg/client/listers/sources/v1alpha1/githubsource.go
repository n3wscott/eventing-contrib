/*
Copyright 2019 The Knative Authors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/knative/eventing-contrib/pkg/apis/sources/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GitHubSourceLister helps list GitHubSources.
type GitHubSourceLister interface {
	// List lists all GitHubSources in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.GitHubSource, err error)
	// GitHubSources returns an object that can list and get GitHubSources.
	GitHubSources(namespace string) GitHubSourceNamespaceLister
	GitHubSourceListerExpansion
}

// gitHubSourceLister implements the GitHubSourceLister interface.
type gitHubSourceLister struct {
	indexer cache.Indexer
}

// NewGitHubSourceLister returns a new GitHubSourceLister.
func NewGitHubSourceLister(indexer cache.Indexer) GitHubSourceLister {
	return &gitHubSourceLister{indexer: indexer}
}

// List lists all GitHubSources in the indexer.
func (s *gitHubSourceLister) List(selector labels.Selector) (ret []*v1alpha1.GitHubSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GitHubSource))
	})
	return ret, err
}

// GitHubSources returns an object that can list and get GitHubSources.
func (s *gitHubSourceLister) GitHubSources(namespace string) GitHubSourceNamespaceLister {
	return gitHubSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GitHubSourceNamespaceLister helps list and get GitHubSources.
type GitHubSourceNamespaceLister interface {
	// List lists all GitHubSources in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.GitHubSource, err error)
	// Get retrieves the GitHubSource from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.GitHubSource, error)
	GitHubSourceNamespaceListerExpansion
}

// gitHubSourceNamespaceLister implements the GitHubSourceNamespaceLister
// interface.
type gitHubSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GitHubSources in the indexer for a given namespace.
func (s gitHubSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GitHubSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GitHubSource))
	})
	return ret, err
}

// Get retrieves the GitHubSource from the indexer for a given namespace and name.
func (s gitHubSourceNamespaceLister) Get(name string) (*v1alpha1.GitHubSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("githubsource"), name)
	}
	return obj.(*v1alpha1.GitHubSource), nil
}
