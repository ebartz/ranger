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
	"github.com/ranger/wrangler/pkg/apply"
	"github.com/ranger/wrangler/pkg/condition"
	"github.com/ranger/wrangler/pkg/generic"
	"github.com/ranger/wrangler/pkg/kv"
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

type MultiClusterAppHandler func(string, *v3.MultiClusterApp) (*v3.MultiClusterApp, error)

type MultiClusterAppController interface {
	generic.ControllerMeta
	MultiClusterAppClient

	OnChange(ctx context.Context, name string, sync MultiClusterAppHandler)
	OnRemove(ctx context.Context, name string, sync MultiClusterAppHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() MultiClusterAppCache
}

type MultiClusterAppClient interface {
	Create(*v3.MultiClusterApp) (*v3.MultiClusterApp, error)
	Update(*v3.MultiClusterApp) (*v3.MultiClusterApp, error)
	UpdateStatus(*v3.MultiClusterApp) (*v3.MultiClusterApp, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.MultiClusterApp, error)
	List(namespace string, opts metav1.ListOptions) (*v3.MultiClusterAppList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.MultiClusterApp, err error)
}

type MultiClusterAppCache interface {
	Get(namespace, name string) (*v3.MultiClusterApp, error)
	List(namespace string, selector labels.Selector) ([]*v3.MultiClusterApp, error)

	AddIndexer(indexName string, indexer MultiClusterAppIndexer)
	GetByIndex(indexName, key string) ([]*v3.MultiClusterApp, error)
}

type MultiClusterAppIndexer func(obj *v3.MultiClusterApp) ([]string, error)

type multiClusterAppController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewMultiClusterAppController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) MultiClusterAppController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &multiClusterAppController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromMultiClusterAppHandlerToHandler(sync MultiClusterAppHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.MultiClusterApp
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.MultiClusterApp))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *multiClusterAppController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.MultiClusterApp))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateMultiClusterAppDeepCopyOnChange(client MultiClusterAppClient, obj *v3.MultiClusterApp, handler func(obj *v3.MultiClusterApp) (*v3.MultiClusterApp, error)) (*v3.MultiClusterApp, error) {
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

func (c *multiClusterAppController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *multiClusterAppController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *multiClusterAppController) OnChange(ctx context.Context, name string, sync MultiClusterAppHandler) {
	c.AddGenericHandler(ctx, name, FromMultiClusterAppHandlerToHandler(sync))
}

func (c *multiClusterAppController) OnRemove(ctx context.Context, name string, sync MultiClusterAppHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromMultiClusterAppHandlerToHandler(sync)))
}

func (c *multiClusterAppController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *multiClusterAppController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *multiClusterAppController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *multiClusterAppController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *multiClusterAppController) Cache() MultiClusterAppCache {
	return &multiClusterAppCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *multiClusterAppController) Create(obj *v3.MultiClusterApp) (*v3.MultiClusterApp, error) {
	result := &v3.MultiClusterApp{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *multiClusterAppController) Update(obj *v3.MultiClusterApp) (*v3.MultiClusterApp, error) {
	result := &v3.MultiClusterApp{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *multiClusterAppController) UpdateStatus(obj *v3.MultiClusterApp) (*v3.MultiClusterApp, error) {
	result := &v3.MultiClusterApp{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *multiClusterAppController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *multiClusterAppController) Get(namespace, name string, options metav1.GetOptions) (*v3.MultiClusterApp, error) {
	result := &v3.MultiClusterApp{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *multiClusterAppController) List(namespace string, opts metav1.ListOptions) (*v3.MultiClusterAppList, error) {
	result := &v3.MultiClusterAppList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *multiClusterAppController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *multiClusterAppController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.MultiClusterApp, error) {
	result := &v3.MultiClusterApp{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type multiClusterAppCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *multiClusterAppCache) Get(namespace, name string) (*v3.MultiClusterApp, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.MultiClusterApp), nil
}

func (c *multiClusterAppCache) List(namespace string, selector labels.Selector) (ret []*v3.MultiClusterApp, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.MultiClusterApp))
	})

	return ret, err
}

func (c *multiClusterAppCache) AddIndexer(indexName string, indexer MultiClusterAppIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.MultiClusterApp))
		},
	}))
}

func (c *multiClusterAppCache) GetByIndex(indexName, key string) (result []*v3.MultiClusterApp, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.MultiClusterApp, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.MultiClusterApp))
	}
	return result, nil
}

type MultiClusterAppStatusHandler func(obj *v3.MultiClusterApp, status v3.MultiClusterAppStatus) (v3.MultiClusterAppStatus, error)

type MultiClusterAppGeneratingHandler func(obj *v3.MultiClusterApp, status v3.MultiClusterAppStatus) ([]runtime.Object, v3.MultiClusterAppStatus, error)

func RegisterMultiClusterAppStatusHandler(ctx context.Context, controller MultiClusterAppController, condition condition.Cond, name string, handler MultiClusterAppStatusHandler) {
	statusHandler := &multiClusterAppStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromMultiClusterAppHandlerToHandler(statusHandler.sync))
}

func RegisterMultiClusterAppGeneratingHandler(ctx context.Context, controller MultiClusterAppController, apply apply.Apply,
	condition condition.Cond, name string, handler MultiClusterAppGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &multiClusterAppGeneratingHandler{
		MultiClusterAppGeneratingHandler: handler,
		apply:                            apply,
		name:                             name,
		gvk:                              controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterMultiClusterAppStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type multiClusterAppStatusHandler struct {
	client    MultiClusterAppClient
	condition condition.Cond
	handler   MultiClusterAppStatusHandler
}

func (a *multiClusterAppStatusHandler) sync(key string, obj *v3.MultiClusterApp) (*v3.MultiClusterApp, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type multiClusterAppGeneratingHandler struct {
	MultiClusterAppGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *multiClusterAppGeneratingHandler) Remove(key string, obj *v3.MultiClusterApp) (*v3.MultiClusterApp, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.MultiClusterApp{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *multiClusterAppGeneratingHandler) Handle(obj *v3.MultiClusterApp, status v3.MultiClusterAppStatus) (v3.MultiClusterAppStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.MultiClusterAppGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
