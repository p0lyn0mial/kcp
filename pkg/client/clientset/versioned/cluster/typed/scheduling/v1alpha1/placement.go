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
	"context"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	schedulingv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/scheduling/v1alpha1"
	schedulingv1alpha1client "github.com/kcp-dev/kcp/pkg/client/clientset/versioned/typed/scheduling/v1alpha1"
)

// PlacementsClusterGetter has a method to return a PlacementClusterInterface.
// A group's cluster client should implement this interface.
type PlacementsClusterGetter interface {
	Placements() PlacementClusterInterface
}

// PlacementClusterInterface can operate on Placements across all clusters,
// or scope down to one cluster and return a schedulingv1alpha1client.PlacementInterface.
type PlacementClusterInterface interface {
	Cluster(logicalcluster.Path) schedulingv1alpha1client.PlacementInterface
	List(ctx context.Context, opts metav1.ListOptions) (*schedulingv1alpha1.PlacementList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type placementsClusterInterface struct {
	clientCache kcpclient.Cache[*schedulingv1alpha1client.SchedulingV1alpha1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *placementsClusterInterface) Cluster(name logicalcluster.Path) schedulingv1alpha1client.PlacementInterface {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(name).Placements()
}

// List returns the entire collection of all Placements across all clusters.
func (c *placementsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*schedulingv1alpha1.PlacementList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).Placements().List(ctx, opts)
}

// Watch begins to watch all Placements across all clusters.
func (c *placementsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).Placements().Watch(ctx, opts)
}
