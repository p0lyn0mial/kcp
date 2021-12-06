/*
Copyright 2021 The KCP Authors.

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
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	v1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apiresource/v1alpha1"
)

// NegotiatedAPIResourceLister helps list NegotiatedAPIResources.
// All objects returned here must be treated as read-only.
type NegotiatedAPIResourceLister interface {
	// List lists all NegotiatedAPIResources in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.NegotiatedAPIResource, err error)
	// ListWithContext lists all NegotiatedAPIResources in the indexer.
	// Objects returned here must be treated as read-only.
	ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1alpha1.NegotiatedAPIResource, err error)
	// Get retrieves the NegotiatedAPIResource from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.NegotiatedAPIResource, error)
	// GetWithContext retrieves the NegotiatedAPIResource from the index for a given name.
	// Objects returned here must be treated as read-only.
	GetWithContext(ctx context.Context, name string) (*v1alpha1.NegotiatedAPIResource, error)
	NegotiatedAPIResourceListerExpansion
}

// negotiatedAPIResourceLister implements the NegotiatedAPIResourceLister interface.
type negotiatedAPIResourceLister struct {
	indexer cache.Indexer
}

// NewNegotiatedAPIResourceLister returns a new NegotiatedAPIResourceLister.
func NewNegotiatedAPIResourceLister(indexer cache.Indexer) NegotiatedAPIResourceLister {
	return &negotiatedAPIResourceLister{indexer: indexer}
}

// List lists all NegotiatedAPIResources in the indexer.
func (s *negotiatedAPIResourceLister) List(selector labels.Selector) (ret []*v1alpha1.NegotiatedAPIResource, err error) {
	return s.ListWithContext(context.Background(), selector)
}

// ListWithContext lists all NegotiatedAPIResources in the indexer.
func (s *negotiatedAPIResourceLister) ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1alpha1.NegotiatedAPIResource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.NegotiatedAPIResource))
	})
	return ret, err
}

// Get retrieves the NegotiatedAPIResource from the index for a given name.
func (s *negotiatedAPIResourceLister) Get(name string) (*v1alpha1.NegotiatedAPIResource, error) {
	return s.GetWithContext(context.Background(), name)
}

// GetWithContext retrieves the NegotiatedAPIResource from the index for a given name.
func (s *negotiatedAPIResourceLister) GetWithContext(ctx context.Context, name string) (*v1alpha1.NegotiatedAPIResource, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("negotiatedapiresource"), name)
	}
	return obj.(*v1alpha1.NegotiatedAPIResource), nil
}
