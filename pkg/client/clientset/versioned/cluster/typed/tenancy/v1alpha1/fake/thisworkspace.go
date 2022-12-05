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

	"github.com/kcp-dev/logicalcluster/v3"

	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	tenancyv1alpha1client "github.com/kcp-dev/kcp/pkg/client/clientset/versioned/typed/tenancy/v1alpha1"
)

var thisWorkspacesResource = schema.GroupVersionResource{Group: "tenancy.kcp.dev", Version: "v1alpha1", Resource: "thisworkspaces"}
var thisWorkspacesKind = schema.GroupVersionKind{Group: "tenancy.kcp.dev", Version: "v1alpha1", Kind: "ThisWorkspace"}

type thisWorkspacesClusterClient struct {
	*kcptesting.Fake
}

// Cluster scopes the client down to a particular cluster.
func (c *thisWorkspacesClusterClient) Cluster(cluster logicalcluster.Path) tenancyv1alpha1client.ThisWorkspaceInterface {
	if cluster == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &thisWorkspacesClient{Fake: c.Fake, Cluster: cluster}
}

// List takes label and field selectors, and returns the list of ThisWorkspaces that match those selectors across all clusters.
func (c *thisWorkspacesClusterClient) List(ctx context.Context, opts metav1.ListOptions) (*tenancyv1alpha1.ThisWorkspaceList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootListAction(thisWorkspacesResource, thisWorkspacesKind, logicalcluster.Wildcard, opts), &tenancyv1alpha1.ThisWorkspaceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &tenancyv1alpha1.ThisWorkspaceList{ListMeta: obj.(*tenancyv1alpha1.ThisWorkspaceList).ListMeta}
	for _, item := range obj.(*tenancyv1alpha1.ThisWorkspaceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested ThisWorkspaces across all clusters.
func (c *thisWorkspacesClusterClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewRootWatchAction(thisWorkspacesResource, logicalcluster.Wildcard, opts))
}

type thisWorkspacesClient struct {
	*kcptesting.Fake
	Cluster logicalcluster.Path
}

func (c *thisWorkspacesClient) Create(ctx context.Context, thisWorkspace *tenancyv1alpha1.ThisWorkspace, opts metav1.CreateOptions) (*tenancyv1alpha1.ThisWorkspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootCreateAction(thisWorkspacesResource, c.Cluster, thisWorkspace), &tenancyv1alpha1.ThisWorkspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.ThisWorkspace), err
}

func (c *thisWorkspacesClient) Update(ctx context.Context, thisWorkspace *tenancyv1alpha1.ThisWorkspace, opts metav1.UpdateOptions) (*tenancyv1alpha1.ThisWorkspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootUpdateAction(thisWorkspacesResource, c.Cluster, thisWorkspace), &tenancyv1alpha1.ThisWorkspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.ThisWorkspace), err
}

func (c *thisWorkspacesClient) UpdateStatus(ctx context.Context, thisWorkspace *tenancyv1alpha1.ThisWorkspace, opts metav1.UpdateOptions) (*tenancyv1alpha1.ThisWorkspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootUpdateSubresourceAction(thisWorkspacesResource, c.Cluster, "status", thisWorkspace), &tenancyv1alpha1.ThisWorkspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.ThisWorkspace), err
}

func (c *thisWorkspacesClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.Invokes(kcptesting.NewRootDeleteActionWithOptions(thisWorkspacesResource, c.Cluster, name, opts), &tenancyv1alpha1.ThisWorkspace{})
	return err
}

func (c *thisWorkspacesClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := kcptesting.NewRootDeleteCollectionAction(thisWorkspacesResource, c.Cluster, listOpts)

	_, err := c.Fake.Invokes(action, &tenancyv1alpha1.ThisWorkspaceList{})
	return err
}

func (c *thisWorkspacesClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*tenancyv1alpha1.ThisWorkspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootGetAction(thisWorkspacesResource, c.Cluster, name), &tenancyv1alpha1.ThisWorkspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.ThisWorkspace), err
}

// List takes label and field selectors, and returns the list of ThisWorkspaces that match those selectors.
func (c *thisWorkspacesClient) List(ctx context.Context, opts metav1.ListOptions) (*tenancyv1alpha1.ThisWorkspaceList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootListAction(thisWorkspacesResource, thisWorkspacesKind, c.Cluster, opts), &tenancyv1alpha1.ThisWorkspaceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &tenancyv1alpha1.ThisWorkspaceList{ListMeta: obj.(*tenancyv1alpha1.ThisWorkspaceList).ListMeta}
	for _, item := range obj.(*tenancyv1alpha1.ThisWorkspaceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *thisWorkspacesClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewRootWatchAction(thisWorkspacesResource, c.Cluster, opts))
}

func (c *thisWorkspacesClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*tenancyv1alpha1.ThisWorkspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootPatchSubresourceAction(thisWorkspacesResource, c.Cluster, name, pt, data, subresources...), &tenancyv1alpha1.ThisWorkspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.ThisWorkspace), err
}
