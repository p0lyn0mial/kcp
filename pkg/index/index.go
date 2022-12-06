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

package index

import (
	"strings"
	"sync"

	"github.com/kcp-dev/logicalcluster/v3"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	tenancyv1beta1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1beta1"
)

// Index implements a mapping from logical cluster to (shard) URL.
type Index interface {
	Lookup(path logicalcluster.Path) (shard string, cluster, canonicalPath logicalcluster.Path, found bool)
	LookupURL(logicalCluster logicalcluster.Path) (url string, canonicalPath logicalcluster.Path, found bool)
}

// PathRewriter can rewrite a logical cluster path before the actual mapping through
// the index data.
type PathRewriter func(segments []string) []string

func New(rewriters []PathRewriter) *State {
	return &State{
		rewriters: rewriters,

		clusterShards:             map[logicalcluster.Path]string{},
		shardWorkspaceNameCluster: map[string]map[logicalcluster.Path]map[string]logicalcluster.Path{},
		shardClusterWorkspaceName: map[string]map[logicalcluster.Path]string{},
		shardClusterParentCluster: map[string]map[logicalcluster.Path]logicalcluster.Path{},
		shardBaseURLs:             map[string]string{},
	}
}

// State watches ClusterWorkspaceShards on the root shard, and then starts informers
// for every ClusterWorkspaceShard, watching the ClusterWorkspaces on them. It then
// updates the workspace index, which maps logical clusters to shard URLs.
type State struct {
	rewriters []PathRewriter

	lock                      sync.RWMutex
	clusterShards             map[logicalcluster.Path]string                                    // logical cluster -> shard name
	shardWorkspaceNameCluster map[string]map[logicalcluster.Path]map[string]logicalcluster.Path // (shard name, logical cluster, workspace name) -> logical cluster
	shardClusterWorkspaceName map[string]map[logicalcluster.Path]string                         // (shard name, logical cluster) -> workspace name
	shardClusterParentCluster map[string]map[logicalcluster.Path]logicalcluster.Path            // (shard name, logical cluster) -> parent logical cluster
	shardBaseURLs             map[string]string                                                 // shard name -> base URL
}

func (c *State) UpsertWorkspace(shard string, ws *tenancyv1beta1.Workspace) {
	if ws.Status.Phase == tenancyv1alpha1.WorkspacePhaseScheduling {
		return
	}
	clusterName := logicalcluster.From(ws)

	c.lock.RLock()
	got := c.shardWorkspaceNameCluster[shard][clusterName.Path()][ws.Name]
	c.lock.RUnlock()

	if got.String() == ws.Status.Cluster {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if got := c.shardWorkspaceNameCluster[shard][clusterName.Path()][ws.Name]; got.String() != ws.Status.Cluster {
		if c.shardWorkspaceNameCluster[shard] == nil {
			c.shardWorkspaceNameCluster[shard] = map[logicalcluster.Path]map[string]logicalcluster.Path{}
			c.shardClusterWorkspaceName[shard] = map[logicalcluster.Path]string{}
			c.shardClusterParentCluster[shard] = map[logicalcluster.Path]logicalcluster.Path{}
		}
		if c.shardWorkspaceNameCluster[shard][clusterName.Path()] == nil {
			c.shardWorkspaceNameCluster[shard][clusterName.Path()] = map[string]logicalcluster.Path{}
		}
		c.shardWorkspaceNameCluster[shard][clusterName.Path()][ws.Name] = logicalcluster.New(ws.Status.Cluster)
		c.shardClusterWorkspaceName[shard][logicalcluster.New(ws.Status.Cluster)] = ws.Name
		c.shardClusterParentCluster[shard][logicalcluster.New(ws.Status.Cluster)] = clusterName.Path()
	}
}

func (c *State) DeleteWorkspace(shard string, ws *tenancyv1beta1.Workspace) {
	if ws.Status.Phase == tenancyv1alpha1.WorkspacePhaseScheduling {
		return
	}
	clusterName := logicalcluster.From(ws)

	c.lock.RLock()
	_, found := c.shardWorkspaceNameCluster[shard][clusterName.Path()][ws.Name]
	c.lock.RUnlock()

	if !found {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if _, found = c.shardWorkspaceNameCluster[shard][clusterName.Path()][ws.Name]; !found {
		return
	}

	delete(c.shardWorkspaceNameCluster[shard][clusterName.Path()], ws.Name)
	if len(c.shardWorkspaceNameCluster[shard][clusterName.Path()]) == 0 {
		delete(c.shardWorkspaceNameCluster[shard], clusterName.Path())
	}
	if len(c.shardWorkspaceNameCluster[shard]) == 0 {
		delete(c.shardWorkspaceNameCluster, shard)
	}

	delete(c.shardClusterWorkspaceName[shard], logicalcluster.New(ws.Status.Cluster))
	if len(c.shardClusterWorkspaceName[shard]) == 0 {
		delete(c.shardClusterWorkspaceName, shard)
	}

	delete(c.shardClusterParentCluster[shard], logicalcluster.New(ws.Status.Cluster))
	if len(c.shardClusterParentCluster[shard]) == 0 {
		delete(c.shardClusterParentCluster, shard)
	}
}

func (c *State) UpsertThisWorkspace(shard string, this *tenancyv1alpha1.ThisWorkspace) {
	clusterName := logicalcluster.From(this)

	c.lock.RLock()
	got := c.clusterShards[clusterName.Path()]
	c.lock.RUnlock()

	if got != shard {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.clusterShards[clusterName.Path()] = shard
	}
}

func (c *State) DeleteThisWorkspace(shard string, this *tenancyv1alpha1.ThisWorkspace) {
	clusterName := logicalcluster.From(this)

	c.lock.RLock()
	got := c.clusterShards[clusterName.Path()]
	c.lock.RUnlock()

	if got == shard {
		c.lock.Lock()
		defer c.lock.Unlock()
		if got := c.clusterShards[clusterName.Path()]; got == shard {
			delete(c.clusterShards, clusterName.Path())
		}
	}
}

func (c *State) UpsertShard(shardName, baseURL string) {
	c.lock.RLock()
	got := c.shardBaseURLs[shardName]
	c.lock.RUnlock()

	if got != baseURL {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.shardBaseURLs[shardName] = baseURL
	}
}

func (c *State) DeleteShard(shardName string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for lc, gotShardName := range c.clusterShards {
		if shardName == gotShardName {
			delete(c.clusterShards, lc)
		}
	}
	delete(c.shardWorkspaceNameCluster, shardName)
	delete(c.shardBaseURLs, shardName)
	delete(c.shardClusterWorkspaceName, shardName)
	delete(c.shardClusterParentCluster, shardName)
}

func (c *State) Lookup(path logicalcluster.Path) (shard string, cluster, canonicalPath logicalcluster.Path, found bool) {
	segments := strings.Split(path.String(), ":")

	for _, rewriter := range c.rewriters {
		segments = rewriter(segments)
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	// walk through index graph to find the final logical cluster and shard
	for i, s := range segments {
		if i == 0 {
			var found bool
			shard, found = c.clusterShards[logicalcluster.New(s)]
			if !found {
				return "", logicalcluster.Path{}, logicalcluster.Path{}, false
			}
			cluster = logicalcluster.New(s)
			continue
		}

		var found bool
		cluster, found = c.shardWorkspaceNameCluster[shard][cluster][s]
		if !found {
			return "", logicalcluster.Path{}, logicalcluster.Path{}, false
		}
		shard, found = c.clusterShards[cluster]
		if !found {
			return "", logicalcluster.Path{}, logicalcluster.Path{}, false
		}
	}

	// walk about the parent graph to reconstruct the canonical workspace path
	var inversePath []string
	ancestor := cluster
	for {
		shard, found = c.clusterShards[ancestor]
		if !found {
			return "", logicalcluster.Path{}, logicalcluster.Path{}, false
		}
		if name, found := c.shardClusterWorkspaceName[shard][ancestor]; !found {
			inversePath = append(inversePath, ancestor.String())
			break
		} else {
			inversePath = append(inversePath, name)
			ancestor = c.shardClusterParentCluster[shard][ancestor]
		}
	}
	for i := len(inversePath) - 1; i >= 0; i-- {
		canonicalPath = canonicalPath.Join(inversePath[i])
	}

	return shard, cluster, canonicalPath, true
}

func (c *State) LookupURL(path logicalcluster.Path) (url string, canonicalPath logicalcluster.Path, found bool) {
	shard, cluster, canonicalPath, found := c.Lookup(path)
	if !found {
		return "", logicalcluster.Path{}, false
	}

	baseURL, found := c.shardBaseURLs[shard]
	if !found {
		return "", logicalcluster.Path{}, false
	}

	return strings.TrimSuffix(baseURL, "/") + cluster.RequestPath(), canonicalPath, true
}
