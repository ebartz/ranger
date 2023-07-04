/*
Copyright 2023 Ranger Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v3

import (
	"context"
	"time"

	"github.com/ranger/lasso/pkg/client"
	"github.com/ranger/lasso/pkg/controller"
	v3 "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ClusterMonitorGraphHandler func(string, *v3.ClusterMonitorGraph) (*v3.ClusterMonitorGraph, error)

type ClusterMonitorGraphController interface {
	generic.ControllerMeta
	ClusterMonitorGraphClient

	OnChange(ctx context.Context, name string, sync ClusterMonitorGraphHandler)
	OnRemove(ctx context.Context, name string, sync ClusterMonitorGraphHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterMonitorGraphCache
}

type ClusterMonitorGraphClient interface {
	Create(*v3.ClusterMonitorGraph) (*v3.ClusterMonitorGraph, error)
	Update(*v3.ClusterMonitorGraph) (*v3.ClusterMonitorGraph, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterMonitorGraph, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ClusterMonitorGraphList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ClusterMonitorGraph, err error)
}

type ClusterMonitorGraphCache interface {
	Get(namespace, name string) (*v3.ClusterMonitorGraph, error)
	List(namespace string, selector labels.Selector) ([]*v3.ClusterMonitorGraph, error)

	AddIndexer(indexName string, indexer ClusterMonitorGraphIndexer)
	GetByIndex(indexName, key string) ([]*v3.ClusterMonitorGraph, error)
}

type ClusterMonitorGraphIndexer func(obj *v3.ClusterMonitorGraph) ([]string, error)

type clusterMonitorGraphController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewClusterMonitorGraphController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ClusterMonitorGraphController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &clusterMonitorGraphController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromClusterMonitorGraphHandlerToHandler(sync ClusterMonitorGraphHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ClusterMonitorGraph
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ClusterMonitorGraph))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *clusterMonitorGraphController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ClusterMonitorGraph))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateClusterMonitorGraphDeepCopyOnChange(client ClusterMonitorGraphClient, obj *v3.ClusterMonitorGraph, handler func(obj *v3.ClusterMonitorGraph) (*v3.ClusterMonitorGraph, error)) (*v3.ClusterMonitorGraph, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *clusterMonitorGraphController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *clusterMonitorGraphController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *clusterMonitorGraphController) OnChange(ctx context.Context, name string, sync ClusterMonitorGraphHandler) {
	c.AddGenericHandler(ctx, name, FromClusterMonitorGraphHandlerToHandler(sync))
}

func (c *clusterMonitorGraphController) OnRemove(ctx context.Context, name string, sync ClusterMonitorGraphHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromClusterMonitorGraphHandlerToHandler(sync)))
}

func (c *clusterMonitorGraphController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *clusterMonitorGraphController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *clusterMonitorGraphController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *clusterMonitorGraphController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *clusterMonitorGraphController) Cache() ClusterMonitorGraphCache {
	return &clusterMonitorGraphCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *clusterMonitorGraphController) Create(obj *v3.ClusterMonitorGraph) (*v3.ClusterMonitorGraph, error) {
	result := &v3.ClusterMonitorGraph{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *clusterMonitorGraphController) Update(obj *v3.ClusterMonitorGraph) (*v3.ClusterMonitorGraph, error) {
	result := &v3.ClusterMonitorGraph{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterMonitorGraphController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *clusterMonitorGraphController) Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterMonitorGraph, error) {
	result := &v3.ClusterMonitorGraph{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *clusterMonitorGraphController) List(namespace string, opts metav1.ListOptions) (*v3.ClusterMonitorGraphList, error) {
	result := &v3.ClusterMonitorGraphList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *clusterMonitorGraphController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *clusterMonitorGraphController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ClusterMonitorGraph, error) {
	result := &v3.ClusterMonitorGraph{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type clusterMonitorGraphCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *clusterMonitorGraphCache) Get(namespace, name string) (*v3.ClusterMonitorGraph, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ClusterMonitorGraph), nil
}

func (c *clusterMonitorGraphCache) List(namespace string, selector labels.Selector) (ret []*v3.ClusterMonitorGraph, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ClusterMonitorGraph))
	})

	return ret, err
}

func (c *clusterMonitorGraphCache) AddIndexer(indexName string, indexer ClusterMonitorGraphIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ClusterMonitorGraph))
		},
	}))
}

func (c *clusterMonitorGraphCache) GetByIndex(indexName, key string) (result []*v3.ClusterMonitorGraph, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ClusterMonitorGraph, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ClusterMonitorGraph))
	}
	return result, nil
}
