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

package v1beta1

import (
	"context"
	"time"

	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	kcpinformers "github.com/kcp-dev/apimachinery/third_party/informers"
	"github.com/kcp-dev/logicalcluster/v3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	tenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
	scopedclientset "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	clientset "github.com/kcp-dev/kcp/pkg/client/clientset/versioned/cluster"
	"github.com/kcp-dev/kcp/pkg/client/informers/externalversions/internalinterfaces"
	tenancyv1beta1listers "github.com/kcp-dev/kcp/pkg/client/listers/tenancy/v1beta1"
)

// WorkspaceClusterInformer provides access to a shared informer and lister for
// Workspaces.
type WorkspaceClusterInformer interface {
	Cluster(logicalcluster.Name) WorkspaceInformer
	Informer() kcpcache.ScopeableSharedIndexInformer
	Lister() tenancyv1beta1listers.WorkspaceClusterLister
}

type workspaceClusterInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewWorkspaceClusterInformer constructs a new informer for Workspace type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewWorkspaceClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredWorkspaceClusterInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredWorkspaceClusterInformer constructs a new informer for Workspace type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredWorkspaceClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) kcpcache.ScopeableSharedIndexInformer {
	return kcpinformers.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TenancyV1beta1().Workspaces().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TenancyV1beta1().Workspaces().Watch(context.TODO(), options)
			},
		},
		&tenancyv1beta1.Workspace{},
		resyncPeriod,
		indexers,
	)
}

func (f *workspaceClusterInformer) defaultInformer(client clientset.ClusterInterface, resyncPeriod time.Duration) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredWorkspaceClusterInformer(client, resyncPeriod, cache.Indexers{
		kcpcache.ClusterIndexName: kcpcache.ClusterIndexFunc,
	},
		f.tweakListOptions,
	)
}

func (f *workspaceClusterInformer) Informer() kcpcache.ScopeableSharedIndexInformer {
	return f.factory.InformerFor(&tenancyv1beta1.Workspace{}, f.defaultInformer)
}

func (f *workspaceClusterInformer) Lister() tenancyv1beta1listers.WorkspaceClusterLister {
	return tenancyv1beta1listers.NewWorkspaceClusterLister(f.Informer().GetIndexer())
}

// WorkspaceInformer provides access to a shared informer and lister for
// Workspaces.
type WorkspaceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() tenancyv1beta1listers.WorkspaceLister
}

func (f *workspaceClusterInformer) Cluster(cluster logicalcluster.Name) WorkspaceInformer {
	return &workspaceInformer{
		informer: f.Informer().Cluster(cluster),
		lister:   f.Lister().Cluster(cluster),
	}
}

type workspaceInformer struct {
	informer cache.SharedIndexInformer
	lister   tenancyv1beta1listers.WorkspaceLister
}

func (f *workspaceInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

func (f *workspaceInformer) Lister() tenancyv1beta1listers.WorkspaceLister {
	return f.lister
}

type workspaceScopedInformer struct {
	factory          internalinterfaces.SharedScopedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func (f *workspaceScopedInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&tenancyv1beta1.Workspace{}, f.defaultInformer)
}

func (f *workspaceScopedInformer) Lister() tenancyv1beta1listers.WorkspaceLister {
	return tenancyv1beta1listers.NewWorkspaceLister(f.Informer().GetIndexer())
}

// NewWorkspaceInformer constructs a new informer for Workspace type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewWorkspaceInformer(client scopedclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredWorkspaceInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredWorkspaceInformer constructs a new informer for Workspace type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredWorkspaceInformer(client scopedclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TenancyV1beta1().Workspaces().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TenancyV1beta1().Workspaces().Watch(context.TODO(), options)
			},
		},
		&tenancyv1beta1.Workspace{},
		resyncPeriod,
		indexers,
	)
}

func (f *workspaceScopedInformer) defaultInformer(client scopedclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredWorkspaceInformer(client, resyncPeriod, cache.Indexers{}, f.tweakListOptions)
}
