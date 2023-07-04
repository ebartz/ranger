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

type CatalogTemplateVersionHandler func(string, *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)

type CatalogTemplateVersionController interface {
	generic.ControllerMeta
	CatalogTemplateVersionClient

	OnChange(ctx context.Context, name string, sync CatalogTemplateVersionHandler)
	OnRemove(ctx context.Context, name string, sync CatalogTemplateVersionHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() CatalogTemplateVersionCache
}

type CatalogTemplateVersionClient interface {
	Create(*v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)
	Update(*v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)
	UpdateStatus(*v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.CatalogTemplateVersion, error)
	List(namespace string, opts metav1.ListOptions) (*v3.CatalogTemplateVersionList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.CatalogTemplateVersion, err error)
}

type CatalogTemplateVersionCache interface {
	Get(namespace, name string) (*v3.CatalogTemplateVersion, error)
	List(namespace string, selector labels.Selector) ([]*v3.CatalogTemplateVersion, error)

	AddIndexer(indexName string, indexer CatalogTemplateVersionIndexer)
	GetByIndex(indexName, key string) ([]*v3.CatalogTemplateVersion, error)
}

type CatalogTemplateVersionIndexer func(obj *v3.CatalogTemplateVersion) ([]string, error)

type catalogTemplateVersionController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewCatalogTemplateVersionController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) CatalogTemplateVersionController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &catalogTemplateVersionController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromCatalogTemplateVersionHandlerToHandler(sync CatalogTemplateVersionHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.CatalogTemplateVersion
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.CatalogTemplateVersion))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *catalogTemplateVersionController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.CatalogTemplateVersion))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateCatalogTemplateVersionDeepCopyOnChange(client CatalogTemplateVersionClient, obj *v3.CatalogTemplateVersion, handler func(obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)) (*v3.CatalogTemplateVersion, error) {
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

func (c *catalogTemplateVersionController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *catalogTemplateVersionController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *catalogTemplateVersionController) OnChange(ctx context.Context, name string, sync CatalogTemplateVersionHandler) {
	c.AddGenericHandler(ctx, name, FromCatalogTemplateVersionHandlerToHandler(sync))
}

func (c *catalogTemplateVersionController) OnRemove(ctx context.Context, name string, sync CatalogTemplateVersionHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromCatalogTemplateVersionHandlerToHandler(sync)))
}

func (c *catalogTemplateVersionController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *catalogTemplateVersionController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *catalogTemplateVersionController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *catalogTemplateVersionController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *catalogTemplateVersionController) Cache() CatalogTemplateVersionCache {
	return &catalogTemplateVersionCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *catalogTemplateVersionController) Create(obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error) {
	result := &v3.CatalogTemplateVersion{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *catalogTemplateVersionController) Update(obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error) {
	result := &v3.CatalogTemplateVersion{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *catalogTemplateVersionController) UpdateStatus(obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error) {
	result := &v3.CatalogTemplateVersion{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *catalogTemplateVersionController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *catalogTemplateVersionController) Get(namespace, name string, options metav1.GetOptions) (*v3.CatalogTemplateVersion, error) {
	result := &v3.CatalogTemplateVersion{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *catalogTemplateVersionController) List(namespace string, opts metav1.ListOptions) (*v3.CatalogTemplateVersionList, error) {
	result := &v3.CatalogTemplateVersionList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *catalogTemplateVersionController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *catalogTemplateVersionController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.CatalogTemplateVersion, error) {
	result := &v3.CatalogTemplateVersion{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type catalogTemplateVersionCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *catalogTemplateVersionCache) Get(namespace, name string) (*v3.CatalogTemplateVersion, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.CatalogTemplateVersion), nil
}

func (c *catalogTemplateVersionCache) List(namespace string, selector labels.Selector) (ret []*v3.CatalogTemplateVersion, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.CatalogTemplateVersion))
	})

	return ret, err
}

func (c *catalogTemplateVersionCache) AddIndexer(indexName string, indexer CatalogTemplateVersionIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.CatalogTemplateVersion))
		},
	}))
}

func (c *catalogTemplateVersionCache) GetByIndex(indexName, key string) (result []*v3.CatalogTemplateVersion, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.CatalogTemplateVersion, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.CatalogTemplateVersion))
	}
	return result, nil
}

type CatalogTemplateVersionStatusHandler func(obj *v3.CatalogTemplateVersion, status v3.TemplateVersionStatus) (v3.TemplateVersionStatus, error)

type CatalogTemplateVersionGeneratingHandler func(obj *v3.CatalogTemplateVersion, status v3.TemplateVersionStatus) ([]runtime.Object, v3.TemplateVersionStatus, error)

func RegisterCatalogTemplateVersionStatusHandler(ctx context.Context, controller CatalogTemplateVersionController, condition condition.Cond, name string, handler CatalogTemplateVersionStatusHandler) {
	statusHandler := &catalogTemplateVersionStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromCatalogTemplateVersionHandlerToHandler(statusHandler.sync))
}

func RegisterCatalogTemplateVersionGeneratingHandler(ctx context.Context, controller CatalogTemplateVersionController, apply apply.Apply,
	condition condition.Cond, name string, handler CatalogTemplateVersionGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &catalogTemplateVersionGeneratingHandler{
		CatalogTemplateVersionGeneratingHandler: handler,
		apply:                                   apply,
		name:                                    name,
		gvk:                                     controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterCatalogTemplateVersionStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type catalogTemplateVersionStatusHandler struct {
	client    CatalogTemplateVersionClient
	condition condition.Cond
	handler   CatalogTemplateVersionStatusHandler
}

func (a *catalogTemplateVersionStatusHandler) sync(key string, obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error) {
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

type catalogTemplateVersionGeneratingHandler struct {
	CatalogTemplateVersionGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *catalogTemplateVersionGeneratingHandler) Remove(key string, obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.CatalogTemplateVersion{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *catalogTemplateVersionGeneratingHandler) Handle(obj *v3.CatalogTemplateVersion, status v3.TemplateVersionStatus) (v3.TemplateVersionStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.CatalogTemplateVersionGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
