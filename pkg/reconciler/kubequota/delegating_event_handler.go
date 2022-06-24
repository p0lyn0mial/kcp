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

package kubequota

import (
	"sync"

	"github.com/kcp-dev/logicalcluster"
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

type delegatingEventHandler struct {
	lock          sync.RWMutex
	eventHandlers map[schema.GroupResource]map[logicalcluster.Name]cache.ResourceEventHandler
}

func (d *delegatingEventHandler) registerEventHandler(
	resource schema.GroupResource,
	informer cache.SharedIndexInformer,
	clusterName logicalcluster.Name,
	h cache.ResourceEventHandler,
) {
	klog.Infof("ANDY registerEventHandler resource=%q, clusterName=%v, grabbing lock", resource, clusterName)
	d.lock.Lock()
	defer func() {
		klog.Infof("ANDY registerEventHandler resource=%q, clusterName=%v, releasing lock", resource, clusterName)
		d.lock.Unlock()
	}()

	groupResourceHandlers, ok := d.eventHandlers[resource]
	if !ok {
		klog.Infof("ANDY registerEventHandler resource=%q, clusterName=%v, need to init for resource", resource, clusterName)
		groupResourceHandlers = map[logicalcluster.Name]cache.ResourceEventHandler{}
		d.eventHandlers[resource] = groupResourceHandlers

		informer.AddEventHandler(
			cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					key, _ := cache.MetaNamespaceKeyFunc(obj)
					h := d.getEventHandler(resource, obj)
					if h == nil {
						klog.Infof("ANDY delegatingEventHandler add event %q %s NO HANDLER", resource, key)
						return
					}
					klog.Infof("ANDY delegatingEventHandler add event %q %s calling handler", resource, key)
					h.OnAdd(obj)
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					key, _ := cache.MetaNamespaceKeyFunc(newObj)
					h := d.getEventHandler(resource, oldObj)
					if h == nil {
						klog.Infof("ANDY delegatingEventHandler update event %q %s NO HANDLER", resource, key)
						return
					}
					klog.Infof("ANDY delegatingEventHandler update event %q %s calling handler", resource, key)
					h.OnUpdate(oldObj, newObj)
				},
				DeleteFunc: func(obj interface{}) {
					key, _ := cache.MetaNamespaceKeyFunc(obj)
					h := d.getEventHandler(resource, obj)
					if h == nil {
						klog.Infof("ANDY delegatingEventHandler delete event %q %s NO HANDLER", resource, key)
						return
					}
					klog.Infof("ANDY delegatingEventHandler delete event %q %s calling handler", resource, key)
					h.OnDelete(obj)
				},
			},
		)
	} else {
		klog.Infof("ANDY registerEventHandler resource=%q, clusterName=%v, no init needed", resource, clusterName)
	}

	klog.Infof("ANDY registerEventHandler resource=%q, clusterName=%v, storing handler", resource, clusterName)
	groupResourceHandlers[clusterName] = h
}

func (d *delegatingEventHandler) getEventHandler(resource schema.GroupResource, obj interface{}) cache.ResourceEventHandler {
	clusterName := clusterNameForObj(obj)
	if clusterName.Empty() {
		return nil
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	groupResourceHandlers, ok := d.eventHandlers[resource]
	if !ok {
		return nil
	}

	return groupResourceHandlers[clusterName]
}
