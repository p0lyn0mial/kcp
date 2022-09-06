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

package server

import (
	"context"
	"fmt"
	"github.com/kcp-dev/logicalcluster/v2"
	"github.com/munnerz/goautoneg"
	kaudit "k8s.io/apiserver/pkg/audit"
	"k8s.io/kubernetes/pkg/genericcontrolplane"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
)

var (
	shardNameRegExp = regexp.MustCompile(`^[a-z0-9-:]{0,61}$`)

	errorScheme = runtime.NewScheme()
	errorCodecs = serializer.NewCodecFactory(errorScheme)
)

func init() {
	errorScheme.AddUnversionedTypes(metav1.Unversioned,
		&metav1.Status{},
	)
}

func WithServiceScope(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if path := req.URL.Path; strings.HasPrefix(path, "/services/cache") {
			path = strings.TrimPrefix(path, "/services/cache")
			req.URL.Path = path
			newURL, err := url.Parse(req.URL.String())
			if err != nil {
				responsewriters.ErrorNegotiated(
					apierrors.NewInternalError(fmt.Errorf("unable to resolve %s, err %w", req.URL.Path, err)),
					errorCodecs, schema.GroupVersion{},
					w, req)
				return
			}
			req.URL = newURL
		}
		handler.ServeHTTP(w, req)
	})
}

// WithShardScope reads a shard name from the URL path and puts it into the context.
// It also trims "/shards/" prefix from the URL.
// If the path doesn't contain the shard name then a default "system:cache:server" name is assigned.
//
// For example:
//
// /shards/*/clusters/*/apis/apis.kcp.dev/v1alpha1/apiexports
//
// /shards/amber/clusters/*/apis/apis.kcp.dev/v1alpha1/apiexports
//
// /shards/sapphire/clusters/system:sapphire/apis/apis.kcp.dev/v1alpha1/apiexports
//
// /shards/amber/clusters/system:amber/apis/apis.kcp.dev/v1alpha1/apiexports
func WithShardScope(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var shardName string
		if path := req.URL.Path; strings.HasPrefix(path, "/shards/") {
			path = strings.TrimPrefix(path, "/shards/")

			i := strings.Index(path, "/")
			if i == -1 {
				responsewriters.ErrorNegotiated(
					apierrors.NewBadRequest(fmt.Sprintf("unable to parse shard: no `/` found in path %s", path)),
					errorCodecs, schema.GroupVersion{},
					w, req)
				return
			}
			shardName, path = path[:i], path[i:]
			req.URL.Path = path
			newURL, err := url.Parse(req.URL.String())
			if err != nil {
				responsewriters.ErrorNegotiated(
					apierrors.NewInternalError(fmt.Errorf("unable to resolve %s, err %w", req.URL.Path, err)),
					errorCodecs, schema.GroupVersion{},
					w, req)
				return
			}
			req.URL = newURL
		}

		var shard request.Shard
		switch {
		case shardName == "*":
			shard = "*"
		case len(shardName) == 0:
			// because we don't store a shard name in an object.
			// requests without a shard name won't be able to find associated data and will fail.
			// as of today we don't instruct controllers used by the apiextention server
			// how to assign/extract a shard name to/from an object.
			// so we need to set a default name here, otherwise these controllers will fail.
			shard = "system:cache:server"
		default:
			if !shardNameRegExp.MatchString(shardName) {
				responsewriters.ErrorNegotiated(
					apierrors.NewBadRequest(fmt.Sprintf("invalid shard: %q does not match the regex", shardName)),
					errorCodecs, schema.GroupVersion{},
					w, req)
				return
			}
			shard = request.Shard(shardName)
		}

		ctx := request.WithShard(req.Context(), shard)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}

var reClusterName = regexp.MustCompile(`^([a-z]([a-z0-9-]{0,61}[a-z0-9])?:)*[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`)

func WithClusterScope(apiHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var clusterName logicalcluster.Name
		if path := req.URL.Path; strings.HasPrefix(path, "/clusters/") {
			path = strings.TrimPrefix(path, "/clusters/")

			i := strings.Index(path, "/")
			if i == -1 {
				responsewriters.ErrorNegotiated(
					apierrors.NewBadRequest(fmt.Sprintf("unable to parse cluster: no `/` found in path %s", path)),
					errorCodecs, schema.GroupVersion{},
					w, req)
				return
			}
			clusterName, path = logicalcluster.New(path[:i]), path[i:]
			req.URL.Path = path
			newURL, err := url.Parse(req.URL.String())
			if err != nil {
				responsewriters.ErrorNegotiated(
					apierrors.NewInternalError(fmt.Errorf("unable to resolve %s, err %w", req.URL.Path, err)),
					errorCodecs, schema.GroupVersion{},
					w, req)
				return
			}
			req.URL = newURL
		} else {
			clusterName = logicalcluster.New(req.Header.Get(logicalcluster.ClusterHeader))
		}

		var cluster request.Cluster

		// This is necessary so wildcard (cross-cluster) partial metadata requests can succeed. The storage layer needs
		// to know if a request is for partial metadata to be able to extract the cluster name from storage keys
		// properly.
		cluster.PartialMetadataRequest = isPartialMetadataRequest(req.Context())

		switch {
		case clusterName == logicalcluster.Wildcard:
			// HACK: just a workaround for testing
			cluster.Wildcard = true
			// fallthrough
			cluster.Name = logicalcluster.Wildcard
		case clusterName.Empty():
			cluster.Name = genericcontrolplane.LocalAdminCluster
		default:
			if !reClusterName.MatchString(clusterName.String()) {
				responsewriters.ErrorNegotiated(
					apierrors.NewBadRequest(fmt.Sprintf("invalid cluster: %q does not match the regex", clusterName)),
					errorCodecs, schema.GroupVersion{},
					w, req)
				return
			}
			cluster.Name = clusterName
		}

		ctx := request.WithCluster(req.Context(), cluster)

		apiHandler.ServeHTTP(w, req.WithContext(ctx))
	}
}

const (
	passthroughHeader   = "X-Kcp-Api-V1-Discovery-Passthrough"
	workspaceAnnotation = "tenancy.kcp.dev/workspace"
)

type (
	acceptHeaderContextKeyType int
	userAgentContextKeyType    int
)

// WithAuditAnnotation initializes audit annotations in the context. Without
// initialization kaudit.AddAuditAnnotation isn't preserved.
func WithAuditAnnotation(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(w, req.WithContext(
			kaudit.WithAuditAnnotations(req.Context()),
		))
	})
}

// WithClusterAnnotation adds the cluster name into the annotation of an audit
// event. Needs initialized annotations.
func WithClusterAnnotation(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cluster := request.ClusterFrom(req.Context())
		if cluster != nil {
			kaudit.AddAuditAnnotation(req.Context(), workspaceAnnotation, cluster.Name.String())
		}

		handler.ServeHTTP(w, req)
	})
}

const (
	// clusterKey is the context key for the request namespace.
	acceptHeaderContextKey acceptHeaderContextKeyType = iota

	// userAgentContextKey is the context key for the request user-agent.
	userAgentContextKey userAgentContextKeyType = iota
)

// WithAcceptHeader makes the Accept header available for code in the handler chain. It is needed for
// Wildcard requests, when finding the CRD with a common schema. For PartialObjectMeta requests we cand
// weaken the schema requirement and allow different schemas across workspaces.
func WithAcceptHeader(apiHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), acceptHeaderContextKey, req.Header.Get("Accept"))
		apiHandler.ServeHTTP(w, req.WithContext(ctx))
	})
}

func WithUserAgent(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), userAgentContextKey, req.Header.Get("User-Agent"))
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}

func isPartialMetadataRequest(ctx context.Context) bool {
	accept := ctx.Value(acceptHeaderContextKey).(string)
	if accept == "" {
		return false
	}

	return isPartialMetadataHeader(accept)
}

func isPartialMetadataHeader(accept string) bool {
	clauses := goautoneg.ParseAccept(accept)
	for _, clause := range clauses {
		if clause.Params["as"] == "PartialObjectMetadata" || clause.Params["as"] == "PartialObjectMetadataList" {
			return true
		}
	}

	return false
}
