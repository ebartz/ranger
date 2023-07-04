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

type MultiClusterAppRevisionHandler func(string, *v3.MultiClusterAppRevision) (*v3.MultiClusterAppRevision, error)

type MultiClusterAppRevisionController interface {
	generic.ControllerMeta
	MultiClusterAppRevisionClient

	OnChange(ctx context.Context, name string, sync MultiClusterAppRevisionHandler)
	OnRemove(ctx context.Context, name string, sync MultiClusterAppRevisionHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() MultiClusterAppRevisionCache
}

type MultiClusterAppRevisionClient interface {
	Create(*v3.MultiClusterAppRevision) (*v3.MultiClusterAppRevision, error)
	Update(*v3.MultiClusterAppRevision) (*v3.MultiClusterAppRevision, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.MultiClusterAppRevision, error)
	List(namespace string, opts metav1.ListOptions) (*v3.MultiClusterAppRevisionList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.MultiClusterAppRevision, err error)
}

type MultiClusterAppRevisionCache interface {
	Get(namespace, name string) (*v3.MultiClusterAppRevision, error)
	List(namespace string, selector labels.Selector) ([]*v3.MultiClusterAppRevision, error)

	AddIndexer(indexName string, indexer MultiClusterAppRevisionIndexer)
	GetByIndex(indexName, key string) ([]*v3.MultiClusterAppRevision, error)
}

type MultiClusterAppRevisionIndexer func(obj *v3.MultiClusterAppRevision) ([]string, error)

type multiClusterAppRevisionController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewMultiClusterAppRevisionController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) MultiClusterAppRevisionController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &multiClusterAppRevisionController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromMultiClusterAppRevisionHandlerToHandler(sync MultiClusterAppRevisionHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.MultiClusterAppRevision
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.MultiClusterAppRevision))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *multiClusterAppRevisionController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.MultiClusterAppRevision))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateMultiClusterAppRevisionDeepCopyOnChange(client MultiClusterAppRevisionClient, obj *v3.MultiClusterAppRevision, handler func(obj *v3.MultiClusterAppRevision) (*v3.MultiClusterAppRevision, error)) (*v3.MultiClusterAppRevision, error) {
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

func (c *multiClusterAppRevisionController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *multiClusterAppRevisionController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *multiClusterAppRevisionController) OnChange(ctx context.Context, name string, sync MultiClusterAppRevisionHandler) {
	c.AddGenericHandler(ctx, name, FromMultiClusterAppRevisionHandlerToHandler(sync))
}

func (c *multiClusterAppRevisionController) OnRemove(ctx context.Context, name string, sync MultiClusterAppRevisionHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromMultiClusterAppRevisionHandlerToHandler(sync)))
}

func (c *multiClusterAppRevisionController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *multiClusterAppRevisionController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *multiClusterAppRevisionController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *multiClusterAppRevisionController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *multiClusterAppRevisionController) Cache() MultiClusterAppRevisionCache {
	return &multiClusterAppRevisionCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *multiClusterAppRevisionController) Create(obj *v3.MultiClusterAppRevision) (*v3.MultiClusterAppRevision, error) {
	result := &v3.MultiClusterAppRevision{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *multiClusterAppRevisionController) Update(obj *v3.MultiClusterAppRevision) (*v3.MultiClusterAppRevision, error) {
	result := &v3.MultiClusterAppRevision{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *multiClusterAppRevisionController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *multiClusterAppRevisionController) Get(namespace, name string, options metav1.GetOptions) (*v3.MultiClusterAppRevision, error) {
	result := &v3.MultiClusterAppRevision{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *multiClusterAppRevisionController) List(namespace string, opts metav1.ListOptions) (*v3.MultiClusterAppRevisionList, error) {
	result := &v3.MultiClusterAppRevisionList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *multiClusterAppRevisionController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *multiClusterAppRevisionController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.MultiClusterAppRevision, error) {
	result := &v3.MultiClusterAppRevision{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type multiClusterAppRevisionCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *multiClusterAppRevisionCache) Get(namespace, name string) (*v3.MultiClusterAppRevision, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.MultiClusterAppRevision), nil
}

func (c *multiClusterAppRevisionCache) List(namespace string, selector labels.Selector) (ret []*v3.MultiClusterAppRevision, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.MultiClusterAppRevision))
	})

	return ret, err
}

func (c *multiClusterAppRevisionCache) AddIndexer(indexName string, indexer MultiClusterAppRevisionIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.MultiClusterAppRevision))
		},
	}))
}

func (c *multiClusterAppRevisionCache) GetByIndex(indexName, key string) (result []*v3.MultiClusterAppRevision, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.MultiClusterAppRevision, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.MultiClusterAppRevision))
	}
	return result, nil
}
