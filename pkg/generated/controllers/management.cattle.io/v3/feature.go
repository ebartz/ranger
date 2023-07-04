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

type FeatureHandler func(string, *v3.Feature) (*v3.Feature, error)

type FeatureController interface {
	generic.ControllerMeta
	FeatureClient

	OnChange(ctx context.Context, name string, sync FeatureHandler)
	OnRemove(ctx context.Context, name string, sync FeatureHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() FeatureCache
}

type FeatureClient interface {
	Create(*v3.Feature) (*v3.Feature, error)
	Update(*v3.Feature) (*v3.Feature, error)
	UpdateStatus(*v3.Feature) (*v3.Feature, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.Feature, error)
	List(opts metav1.ListOptions) (*v3.FeatureList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.Feature, err error)
}

type FeatureCache interface {
	Get(name string) (*v3.Feature, error)
	List(selector labels.Selector) ([]*v3.Feature, error)

	AddIndexer(indexName string, indexer FeatureIndexer)
	GetByIndex(indexName, key string) ([]*v3.Feature, error)
}

type FeatureIndexer func(obj *v3.Feature) ([]string, error)

type featureController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewFeatureController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) FeatureController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &featureController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromFeatureHandlerToHandler(sync FeatureHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.Feature
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.Feature))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *featureController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.Feature))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateFeatureDeepCopyOnChange(client FeatureClient, obj *v3.Feature, handler func(obj *v3.Feature) (*v3.Feature, error)) (*v3.Feature, error) {
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

func (c *featureController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *featureController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *featureController) OnChange(ctx context.Context, name string, sync FeatureHandler) {
	c.AddGenericHandler(ctx, name, FromFeatureHandlerToHandler(sync))
}

func (c *featureController) OnRemove(ctx context.Context, name string, sync FeatureHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromFeatureHandlerToHandler(sync)))
}

func (c *featureController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *featureController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *featureController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *featureController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *featureController) Cache() FeatureCache {
	return &featureCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *featureController) Create(obj *v3.Feature) (*v3.Feature, error) {
	result := &v3.Feature{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *featureController) Update(obj *v3.Feature) (*v3.Feature, error) {
	result := &v3.Feature{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *featureController) UpdateStatus(obj *v3.Feature) (*v3.Feature, error) {
	result := &v3.Feature{}
	return result, c.client.UpdateStatus(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *featureController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *featureController) Get(name string, options metav1.GetOptions) (*v3.Feature, error) {
	result := &v3.Feature{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *featureController) List(opts metav1.ListOptions) (*v3.FeatureList, error) {
	result := &v3.FeatureList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *featureController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *featureController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v3.Feature, error) {
	result := &v3.Feature{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type featureCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *featureCache) Get(name string) (*v3.Feature, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.Feature), nil
}

func (c *featureCache) List(selector labels.Selector) (ret []*v3.Feature, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.Feature))
	})

	return ret, err
}

func (c *featureCache) AddIndexer(indexName string, indexer FeatureIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.Feature))
		},
	}))
}

func (c *featureCache) GetByIndex(indexName, key string) (result []*v3.Feature, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.Feature, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.Feature))
	}
	return result, nil
}

type FeatureStatusHandler func(obj *v3.Feature, status v3.FeatureStatus) (v3.FeatureStatus, error)

type FeatureGeneratingHandler func(obj *v3.Feature, status v3.FeatureStatus) ([]runtime.Object, v3.FeatureStatus, error)

func RegisterFeatureStatusHandler(ctx context.Context, controller FeatureController, condition condition.Cond, name string, handler FeatureStatusHandler) {
	statusHandler := &featureStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromFeatureHandlerToHandler(statusHandler.sync))
}

func RegisterFeatureGeneratingHandler(ctx context.Context, controller FeatureController, apply apply.Apply,
	condition condition.Cond, name string, handler FeatureGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &featureGeneratingHandler{
		FeatureGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterFeatureStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type featureStatusHandler struct {
	client    FeatureClient
	condition condition.Cond
	handler   FeatureStatusHandler
}

func (a *featureStatusHandler) sync(key string, obj *v3.Feature) (*v3.Feature, error) {
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

type featureGeneratingHandler struct {
	FeatureGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *featureGeneratingHandler) Remove(key string, obj *v3.Feature) (*v3.Feature, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.Feature{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *featureGeneratingHandler) Handle(obj *v3.Feature, status v3.FeatureStatus) (v3.FeatureStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.FeatureGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
