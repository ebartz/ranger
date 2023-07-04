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
	v3 "github.com/ranger/ranger/pkg/apis/project.cattle.io/v3"
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

type WorkloadHandler func(string, *v3.Workload) (*v3.Workload, error)

type WorkloadController interface {
	generic.ControllerMeta
	WorkloadClient

	OnChange(ctx context.Context, name string, sync WorkloadHandler)
	OnRemove(ctx context.Context, name string, sync WorkloadHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() WorkloadCache
}

type WorkloadClient interface {
	Create(*v3.Workload) (*v3.Workload, error)
	Update(*v3.Workload) (*v3.Workload, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.Workload, error)
	List(namespace string, opts metav1.ListOptions) (*v3.WorkloadList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.Workload, err error)
}

type WorkloadCache interface {
	Get(namespace, name string) (*v3.Workload, error)
	List(namespace string, selector labels.Selector) ([]*v3.Workload, error)

	AddIndexer(indexName string, indexer WorkloadIndexer)
	GetByIndex(indexName, key string) ([]*v3.Workload, error)
}

type WorkloadIndexer func(obj *v3.Workload) ([]string, error)

type workloadController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewWorkloadController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) WorkloadController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &workloadController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromWorkloadHandlerToHandler(sync WorkloadHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.Workload
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.Workload))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *workloadController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.Workload))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateWorkloadDeepCopyOnChange(client WorkloadClient, obj *v3.Workload, handler func(obj *v3.Workload) (*v3.Workload, error)) (*v3.Workload, error) {
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

func (c *workloadController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *workloadController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *workloadController) OnChange(ctx context.Context, name string, sync WorkloadHandler) {
	c.AddGenericHandler(ctx, name, FromWorkloadHandlerToHandler(sync))
}

func (c *workloadController) OnRemove(ctx context.Context, name string, sync WorkloadHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromWorkloadHandlerToHandler(sync)))
}

func (c *workloadController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *workloadController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *workloadController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *workloadController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *workloadController) Cache() WorkloadCache {
	return &workloadCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *workloadController) Create(obj *v3.Workload) (*v3.Workload, error) {
	result := &v3.Workload{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *workloadController) Update(obj *v3.Workload) (*v3.Workload, error) {
	result := &v3.Workload{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *workloadController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *workloadController) Get(namespace, name string, options metav1.GetOptions) (*v3.Workload, error) {
	result := &v3.Workload{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *workloadController) List(namespace string, opts metav1.ListOptions) (*v3.WorkloadList, error) {
	result := &v3.WorkloadList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *workloadController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *workloadController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.Workload, error) {
	result := &v3.Workload{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type workloadCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *workloadCache) Get(namespace, name string) (*v3.Workload, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.Workload), nil
}

func (c *workloadCache) List(namespace string, selector labels.Selector) (ret []*v3.Workload, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.Workload))
	})

	return ret, err
}

func (c *workloadCache) AddIndexer(indexName string, indexer WorkloadIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.Workload))
		},
	}))
}

func (c *workloadCache) GetByIndex(indexName, key string) (result []*v3.Workload, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.Workload, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.Workload))
	}
	return result, nil
}
