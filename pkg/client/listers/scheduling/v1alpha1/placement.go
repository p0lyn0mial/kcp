//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The KCP Authors.

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

// Code generated by kcp code-generator. DO NOT EDIT.

package v1alpha1

import (
	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	"github.com/kcp-dev/logicalcluster/v3"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	schedulingv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/scheduling/v1alpha1"
)

// PlacementClusterLister can list Placements across all workspaces, or scope down to a PlacementLister for one workspace.
// All objects returned here must be treated as read-only.
type PlacementClusterLister interface {
	// List lists all Placements in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*schedulingv1alpha1.Placement, err error)
	// Cluster returns a lister that can list and get Placements in one workspace.
	Cluster(cluster logicalcluster.Name) PlacementLister
	PlacementClusterListerExpansion
}

type placementClusterLister struct {
	indexer cache.Indexer
}

// NewPlacementClusterLister returns a new PlacementClusterLister.
// We assume that the indexer:
// - is fed by a cross-workspace LIST+WATCH
// - uses kcpcache.MetaClusterNamespaceKeyFunc as the key function
// - has the kcpcache.ClusterIndex as an index
func NewPlacementClusterLister(indexer cache.Indexer) *placementClusterLister {
	return &placementClusterLister{indexer: indexer}
}

// List lists all Placements in the indexer across all workspaces.
func (s *placementClusterLister) List(selector labels.Selector) (ret []*schedulingv1alpha1.Placement, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*schedulingv1alpha1.Placement))
	})
	return ret, err
}

// Cluster scopes the lister to one workspace, allowing users to list and get Placements.
func (s *placementClusterLister) Cluster(cluster logicalcluster.Name) PlacementLister {
	return &placementLister{indexer: s.indexer, cluster: cluster}
}

// PlacementLister can list all Placements, or get one in particular.
// All objects returned here must be treated as read-only.
type PlacementLister interface {
	// List lists all Placements in the workspace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*schedulingv1alpha1.Placement, err error)
	// Get retrieves the Placement from the indexer for a given workspace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*schedulingv1alpha1.Placement, error)
	PlacementListerExpansion
}

// placementLister can list all Placements inside a workspace.
type placementLister struct {
	indexer cache.Indexer
	cluster logicalcluster.Name
}

// List lists all Placements in the indexer for a workspace.
func (s *placementLister) List(selector labels.Selector) (ret []*schedulingv1alpha1.Placement, err error) {
	err = kcpcache.ListAllByCluster(s.indexer, s.cluster, selector, func(i interface{}) {
		ret = append(ret, i.(*schedulingv1alpha1.Placement))
	})
	return ret, err
}

// Get retrieves the Placement from the indexer for a given workspace and name.
func (s *placementLister) Get(name string) (*schedulingv1alpha1.Placement, error) {
	key := kcpcache.ToClusterAwareKey(s.cluster.String(), "", name)
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schedulingv1alpha1.Resource("Placement"), name)
	}
	return obj.(*schedulingv1alpha1.Placement), nil
}

// NewPlacementLister returns a new PlacementLister.
// We assume that the indexer:
// - is fed by a workspace-scoped LIST+WATCH
// - uses cache.MetaNamespaceKeyFunc as the key function
func NewPlacementLister(indexer cache.Indexer) *placementScopedLister {
	return &placementScopedLister{indexer: indexer}
}

// placementScopedLister can list all Placements inside a workspace.
type placementScopedLister struct {
	indexer cache.Indexer
}

// List lists all Placements in the indexer for a workspace.
func (s *placementScopedLister) List(selector labels.Selector) (ret []*schedulingv1alpha1.Placement, err error) {
	err = cache.ListAll(s.indexer, selector, func(i interface{}) {
		ret = append(ret, i.(*schedulingv1alpha1.Placement))
	})
	return ret, err
}

// Get retrieves the Placement from the indexer for a given workspace and name.
func (s *placementScopedLister) Get(name string) (*schedulingv1alpha1.Placement, error) {
	key := name
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schedulingv1alpha1.Resource("Placement"), name)
	}
	return obj.(*schedulingv1alpha1.Placement), nil
}
