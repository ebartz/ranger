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

type ClusterRoleTemplateBindingHandler func(string, *v3.ClusterRoleTemplateBinding) (*v3.ClusterRoleTemplateBinding, error)

type ClusterRoleTemplateBindingController interface {
	generic.ControllerMeta
	ClusterRoleTemplateBindingClient

	OnChange(ctx context.Context, name string, sync ClusterRoleTemplateBindingHandler)
	OnRemove(ctx context.Context, name string, sync ClusterRoleTemplateBindingHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterRoleTemplateBindingCache
}

type ClusterRoleTemplateBindingClient interface {
	Create(*v3.ClusterRoleTemplateBinding) (*v3.ClusterRoleTemplateBinding, error)
	Update(*v3.ClusterRoleTemplateBinding) (*v3.ClusterRoleTemplateBinding, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterRoleTemplateBinding, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ClusterRoleTemplateBindingList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ClusterRoleTemplateBinding, err error)
}

type ClusterRoleTemplateBindingCache interface {
	Get(namespace, name string) (*v3.ClusterRoleTemplateBinding, error)
	List(namespace string, selector labels.Selector) ([]*v3.ClusterRoleTemplateBinding, error)

	AddIndexer(indexName string, indexer ClusterRoleTemplateBindingIndexer)
	GetByIndex(indexName, key string) ([]*v3.ClusterRoleTemplateBinding, error)
}

type ClusterRoleTemplateBindingIndexer func(obj *v3.ClusterRoleTemplateBinding) ([]string, error)

type clusterRoleTemplateBindingController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewClusterRoleTemplateBindingController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ClusterRoleTemplateBindingController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &clusterRoleTemplateBindingController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromClusterRoleTemplateBindingHandlerToHandler(sync ClusterRoleTemplateBindingHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ClusterRoleTemplateBinding
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ClusterRoleTemplateBinding))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *clusterRoleTemplateBindingController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ClusterRoleTemplateBinding))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateClusterRoleTemplateBindingDeepCopyOnChange(client ClusterRoleTemplateBindingClient, obj *v3.ClusterRoleTemplateBinding, handler func(obj *v3.ClusterRoleTemplateBinding) (*v3.ClusterRoleTemplateBinding, error)) (*v3.ClusterRoleTemplateBinding, error) {
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

func (c *clusterRoleTemplateBindingController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *clusterRoleTemplateBindingController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *clusterRoleTemplateBindingController) OnChange(ctx context.Context, name string, sync ClusterRoleTemplateBindingHandler) {
	c.AddGenericHandler(ctx, name, FromClusterRoleTemplateBindingHandlerToHandler(sync))
}

func (c *clusterRoleTemplateBindingController) OnRemove(ctx context.Context, name string, sync ClusterRoleTemplateBindingHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromClusterRoleTemplateBindingHandlerToHandler(sync)))
}

func (c *clusterRoleTemplateBindingController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *clusterRoleTemplateBindingController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *clusterRoleTemplateBindingController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *clusterRoleTemplateBindingController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *clusterRoleTemplateBindingController) Cache() ClusterRoleTemplateBindingCache {
	return &clusterRoleTemplateBindingCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *clusterRoleTemplateBindingController) Create(obj *v3.ClusterRoleTemplateBinding) (*v3.ClusterRoleTemplateBinding, error) {
	result := &v3.ClusterRoleTemplateBinding{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *clusterRoleTemplateBindingController) Update(obj *v3.ClusterRoleTemplateBinding) (*v3.ClusterRoleTemplateBinding, error) {
	result := &v3.ClusterRoleTemplateBinding{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterRoleTemplateBindingController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *clusterRoleTemplateBindingController) Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterRoleTemplateBinding, error) {
	result := &v3.ClusterRoleTemplateBinding{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *clusterRoleTemplateBindingController) List(namespace string, opts metav1.ListOptions) (*v3.ClusterRoleTemplateBindingList, error) {
	result := &v3.ClusterRoleTemplateBindingList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *clusterRoleTemplateBindingController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *clusterRoleTemplateBindingController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ClusterRoleTemplateBinding, error) {
	result := &v3.ClusterRoleTemplateBinding{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type clusterRoleTemplateBindingCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *clusterRoleTemplateBindingCache) Get(namespace, name string) (*v3.ClusterRoleTemplateBinding, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ClusterRoleTemplateBinding), nil
}

func (c *clusterRoleTemplateBindingCache) List(namespace string, selector labels.Selector) (ret []*v3.ClusterRoleTemplateBinding, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ClusterRoleTemplateBinding))
	})

	return ret, err
}

func (c *clusterRoleTemplateBindingCache) AddIndexer(indexName string, indexer ClusterRoleTemplateBindingIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ClusterRoleTemplateBinding))
		},
	}))
}

func (c *clusterRoleTemplateBindingCache) GetByIndex(indexName, key string) (result []*v3.ClusterRoleTemplateBinding, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ClusterRoleTemplateBinding, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ClusterRoleTemplateBinding))
	}
	return result, nil
}
