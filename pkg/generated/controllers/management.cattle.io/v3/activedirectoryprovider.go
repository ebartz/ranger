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

type ActiveDirectoryProviderHandler func(string, *v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error)

type ActiveDirectoryProviderController interface {
	generic.ControllerMeta
	ActiveDirectoryProviderClient

	OnChange(ctx context.Context, name string, sync ActiveDirectoryProviderHandler)
	OnRemove(ctx context.Context, name string, sync ActiveDirectoryProviderHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() ActiveDirectoryProviderCache
}

type ActiveDirectoryProviderClient interface {
	Create(*v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error)
	Update(*v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.ActiveDirectoryProvider, error)
	List(opts metav1.ListOptions) (*v3.ActiveDirectoryProviderList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ActiveDirectoryProvider, err error)
}

type ActiveDirectoryProviderCache interface {
	Get(name string) (*v3.ActiveDirectoryProvider, error)
	List(selector labels.Selector) ([]*v3.ActiveDirectoryProvider, error)

	AddIndexer(indexName string, indexer ActiveDirectoryProviderIndexer)
	GetByIndex(indexName, key string) ([]*v3.ActiveDirectoryProvider, error)
}

type ActiveDirectoryProviderIndexer func(obj *v3.ActiveDirectoryProvider) ([]string, error)

type activeDirectoryProviderController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewActiveDirectoryProviderController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ActiveDirectoryProviderController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &activeDirectoryProviderController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromActiveDirectoryProviderHandlerToHandler(sync ActiveDirectoryProviderHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ActiveDirectoryProvider
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ActiveDirectoryProvider))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *activeDirectoryProviderController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ActiveDirectoryProvider))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateActiveDirectoryProviderDeepCopyOnChange(client ActiveDirectoryProviderClient, obj *v3.ActiveDirectoryProvider, handler func(obj *v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error)) (*v3.ActiveDirectoryProvider, error) {
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

func (c *activeDirectoryProviderController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *activeDirectoryProviderController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *activeDirectoryProviderController) OnChange(ctx context.Context, name string, sync ActiveDirectoryProviderHandler) {
	c.AddGenericHandler(ctx, name, FromActiveDirectoryProviderHandlerToHandler(sync))
}

func (c *activeDirectoryProviderController) OnRemove(ctx context.Context, name string, sync ActiveDirectoryProviderHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromActiveDirectoryProviderHandlerToHandler(sync)))
}

func (c *activeDirectoryProviderController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *activeDirectoryProviderController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *activeDirectoryProviderController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *activeDirectoryProviderController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *activeDirectoryProviderController) Cache() ActiveDirectoryProviderCache {
	return &activeDirectoryProviderCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *activeDirectoryProviderController) Create(obj *v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error) {
	result := &v3.ActiveDirectoryProvider{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *activeDirectoryProviderController) Update(obj *v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error) {
	result := &v3.ActiveDirectoryProvider{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *activeDirectoryProviderController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *activeDirectoryProviderController) Get(name string, options metav1.GetOptions) (*v3.ActiveDirectoryProvider, error) {
	result := &v3.ActiveDirectoryProvider{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *activeDirectoryProviderController) List(opts metav1.ListOptions) (*v3.ActiveDirectoryProviderList, error) {
	result := &v3.ActiveDirectoryProviderList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *activeDirectoryProviderController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *activeDirectoryProviderController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ActiveDirectoryProvider, error) {
	result := &v3.ActiveDirectoryProvider{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type activeDirectoryProviderCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *activeDirectoryProviderCache) Get(name string) (*v3.ActiveDirectoryProvider, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ActiveDirectoryProvider), nil
}

func (c *activeDirectoryProviderCache) List(selector labels.Selector) (ret []*v3.ActiveDirectoryProvider, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ActiveDirectoryProvider))
	})

	return ret, err
}

func (c *activeDirectoryProviderCache) AddIndexer(indexName string, indexer ActiveDirectoryProviderIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ActiveDirectoryProvider))
		},
	}))
}

func (c *activeDirectoryProviderCache) GetByIndex(indexName, key string) (result []*v3.ActiveDirectoryProvider, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ActiveDirectoryProvider, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ActiveDirectoryProvider))
	}
	return result, nil
}
