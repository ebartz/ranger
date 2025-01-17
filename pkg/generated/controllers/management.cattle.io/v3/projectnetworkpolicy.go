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

type ProjectNetworkPolicyHandler func(string, *v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error)

type ProjectNetworkPolicyController interface {
	generic.ControllerMeta
	ProjectNetworkPolicyClient

	OnChange(ctx context.Context, name string, sync ProjectNetworkPolicyHandler)
	OnRemove(ctx context.Context, name string, sync ProjectNetworkPolicyHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ProjectNetworkPolicyCache
}

type ProjectNetworkPolicyClient interface {
	Create(*v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error)
	Update(*v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ProjectNetworkPolicy, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ProjectNetworkPolicyList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ProjectNetworkPolicy, err error)
}

type ProjectNetworkPolicyCache interface {
	Get(namespace, name string) (*v3.ProjectNetworkPolicy, error)
	List(namespace string, selector labels.Selector) ([]*v3.ProjectNetworkPolicy, error)

	AddIndexer(indexName string, indexer ProjectNetworkPolicyIndexer)
	GetByIndex(indexName, key string) ([]*v3.ProjectNetworkPolicy, error)
}

type ProjectNetworkPolicyIndexer func(obj *v3.ProjectNetworkPolicy) ([]string, error)

type projectNetworkPolicyController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewProjectNetworkPolicyController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ProjectNetworkPolicyController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &projectNetworkPolicyController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromProjectNetworkPolicyHandlerToHandler(sync ProjectNetworkPolicyHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ProjectNetworkPolicy
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ProjectNetworkPolicy))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *projectNetworkPolicyController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ProjectNetworkPolicy))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateProjectNetworkPolicyDeepCopyOnChange(client ProjectNetworkPolicyClient, obj *v3.ProjectNetworkPolicy, handler func(obj *v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error)) (*v3.ProjectNetworkPolicy, error) {
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

func (c *projectNetworkPolicyController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *projectNetworkPolicyController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *projectNetworkPolicyController) OnChange(ctx context.Context, name string, sync ProjectNetworkPolicyHandler) {
	c.AddGenericHandler(ctx, name, FromProjectNetworkPolicyHandlerToHandler(sync))
}

func (c *projectNetworkPolicyController) OnRemove(ctx context.Context, name string, sync ProjectNetworkPolicyHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromProjectNetworkPolicyHandlerToHandler(sync)))
}

func (c *projectNetworkPolicyController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *projectNetworkPolicyController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *projectNetworkPolicyController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *projectNetworkPolicyController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *projectNetworkPolicyController) Cache() ProjectNetworkPolicyCache {
	return &projectNetworkPolicyCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *projectNetworkPolicyController) Create(obj *v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error) {
	result := &v3.ProjectNetworkPolicy{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *projectNetworkPolicyController) Update(obj *v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error) {
	result := &v3.ProjectNetworkPolicy{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *projectNetworkPolicyController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *projectNetworkPolicyController) Get(namespace, name string, options metav1.GetOptions) (*v3.ProjectNetworkPolicy, error) {
	result := &v3.ProjectNetworkPolicy{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *projectNetworkPolicyController) List(namespace string, opts metav1.ListOptions) (*v3.ProjectNetworkPolicyList, error) {
	result := &v3.ProjectNetworkPolicyList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *projectNetworkPolicyController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *projectNetworkPolicyController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ProjectNetworkPolicy, error) {
	result := &v3.ProjectNetworkPolicy{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type projectNetworkPolicyCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *projectNetworkPolicyCache) Get(namespace, name string) (*v3.ProjectNetworkPolicy, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ProjectNetworkPolicy), nil
}

func (c *projectNetworkPolicyCache) List(namespace string, selector labels.Selector) (ret []*v3.ProjectNetworkPolicy, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ProjectNetworkPolicy))
	})

	return ret, err
}

func (c *projectNetworkPolicyCache) AddIndexer(indexName string, indexer ProjectNetworkPolicyIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ProjectNetworkPolicy))
		},
	}))
}

func (c *projectNetworkPolicyCache) GetByIndex(indexName, key string) (result []*v3.ProjectNetworkPolicy, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ProjectNetworkPolicy, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ProjectNetworkPolicy))
	}
	return result, nil
}
