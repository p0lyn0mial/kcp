/*
Copyright 2022 The KCP Authors.

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

package tenancy

import (
	"strings"

	"github.com/kcp-dev/logicalcluster/v3"

	"k8s.io/klog/v2"
)

const (
	// LogicalClusterPathAnnotationKey is the annotation key for the logical cluster path
	// put on objects that are referenced by path by other objects.
	//
	// If this annotation exists, the system will maintain the annotation value.
	LogicalClusterPathAnnotationKey = "tenancy.kcp.dev/path"
)

// Object is a local interface representation of the Kubernetes metav1.Object, to avoid dependencies on
// k8s.io/apimachinery.
type Object interface {
	GetAnnotations() map[string]string
}

// Cluster is a logical cluster name, not a workspace path.
//
// +kubebuilder:validation:Pattern:="^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
type Cluster string

func (c Cluster) Path() logicalcluster.Name {
	return logicalcluster.New(string(c))
}

func (c Cluster) String() string {
	return string(c)
}

func From(obj Object) Cluster {
	return Cluster(logicalcluster.From(obj).String())
}

// TemporaryCanonicalPath maps a cluster name to the canonical workspace path
// for that cluster. This is temporary, and it will be replaced by some cached
// mapping backed by the workspace index, probably of the front-proxy.
func TemporaryCanonicalPath(c Cluster) logicalcluster.Name {
	path := logicalcluster.New(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(c), "-", "§"), "-", ":"), "§", "-"))

	logger := klog.Background()
	logger.V(1).Info("TemporaryCanonicalPath", "cluster", c, "path", path) // intentionally noisy output

	return path
}

// TemporaryClusterFrom returns the cluster name for a given workspace path.
// This is temporary, and it will be replaced by some cached mapping backed
// by the workspace index, probably of the front-proxy.
func TemporaryClusterFrom(path logicalcluster.Name) Cluster {
	cluster := Cluster(strings.ReplaceAll(strings.ReplaceAll(path.String(), "-", "--"), ":", "-"))

	logger := klog.Background()
	logger.V(1).Info("TemporaryClusterFrom", "path", path, "cluster", cluster) // intentionally noisy output

	return cluster
}
