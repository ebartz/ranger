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

type ClusterTemplateRevisionHandler func(string, *v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error)

type ClusterTemplateRevisionController interface {
	generic.ControllerMeta
	ClusterTemplateRevisionClient

	OnChange(ctx context.Context, name string, sync ClusterTemplateRevisionHandler)
	OnRemove(ctx context.Context, name string, sync ClusterTemplateRevisionHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterTemplateRevisionCache
}

type ClusterTemplateRevisionClient interface {
	Create(*v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error)
	Update(*v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error)
	UpdateStatus(*v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterTemplateRevision, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ClusterTemplateRevisionList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ClusterTemplateRevision, err error)
}

type ClusterTemplateRevisionCache interface {
	Get(namespace, name string) (*v3.ClusterTemplateRevision, error)
	List(namespace string, selector labels.Selector) ([]*v3.ClusterTemplateRevision, error)

	AddIndexer(indexName string, indexer ClusterTemplateRevisionIndexer)
	GetByIndex(indexName, key string) ([]*v3.ClusterTemplateRevision, error)
}

type ClusterTemplateRevisionIndexer func(obj *v3.ClusterTemplateRevision) ([]string, error)

type clusterTemplateRevisionController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewClusterTemplateRevisionController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ClusterTemplateRevisionController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &clusterTemplateRevisionController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromClusterTemplateRevisionHandlerToHandler(sync ClusterTemplateRevisionHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ClusterTemplateRevision
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ClusterTemplateRevision))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *clusterTemplateRevisionController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ClusterTemplateRevision))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateClusterTemplateRevisionDeepCopyOnChange(client ClusterTemplateRevisionClient, obj *v3.ClusterTemplateRevision, handler func(obj *v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error)) (*v3.ClusterTemplateRevision, error) {
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

func (c *clusterTemplateRevisionController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *clusterTemplateRevisionController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *clusterTemplateRevisionController) OnChange(ctx context.Context, name string, sync ClusterTemplateRevisionHandler) {
	c.AddGenericHandler(ctx, name, FromClusterTemplateRevisionHandlerToHandler(sync))
}

func (c *clusterTemplateRevisionController) OnRemove(ctx context.Context, name string, sync ClusterTemplateRevisionHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromClusterTemplateRevisionHandlerToHandler(sync)))
}

func (c *clusterTemplateRevisionController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *clusterTemplateRevisionController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *clusterTemplateRevisionController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *clusterTemplateRevisionController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *clusterTemplateRevisionController) Cache() ClusterTemplateRevisionCache {
	return &clusterTemplateRevisionCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *clusterTemplateRevisionController) Create(obj *v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error) {
	result := &v3.ClusterTemplateRevision{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *clusterTemplateRevisionController) Update(obj *v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error) {
	result := &v3.ClusterTemplateRevision{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterTemplateRevisionController) UpdateStatus(obj *v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error) {
	result := &v3.ClusterTemplateRevision{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterTemplateRevisionController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *clusterTemplateRevisionController) Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterTemplateRevision, error) {
	result := &v3.ClusterTemplateRevision{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *clusterTemplateRevisionController) List(namespace string, opts metav1.ListOptions) (*v3.ClusterTemplateRevisionList, error) {
	result := &v3.ClusterTemplateRevisionList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *clusterTemplateRevisionController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *clusterTemplateRevisionController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ClusterTemplateRevision, error) {
	result := &v3.ClusterTemplateRevision{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type clusterTemplateRevisionCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *clusterTemplateRevisionCache) Get(namespace, name string) (*v3.ClusterTemplateRevision, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ClusterTemplateRevision), nil
}

func (c *clusterTemplateRevisionCache) List(namespace string, selector labels.Selector) (ret []*v3.ClusterTemplateRevision, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ClusterTemplateRevision))
	})

	return ret, err
}

func (c *clusterTemplateRevisionCache) AddIndexer(indexName string, indexer ClusterTemplateRevisionIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ClusterTemplateRevision))
		},
	}))
}

func (c *clusterTemplateRevisionCache) GetByIndex(indexName, key string) (result []*v3.ClusterTemplateRevision, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ClusterTemplateRevision, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ClusterTemplateRevision))
	}
	return result, nil
}

type ClusterTemplateRevisionStatusHandler func(obj *v3.ClusterTemplateRevision, status v3.ClusterTemplateRevisionStatus) (v3.ClusterTemplateRevisionStatus, error)

type ClusterTemplateRevisionGeneratingHandler func(obj *v3.ClusterTemplateRevision, status v3.ClusterTemplateRevisionStatus) ([]runtime.Object, v3.ClusterTemplateRevisionStatus, error)

func RegisterClusterTemplateRevisionStatusHandler(ctx context.Context, controller ClusterTemplateRevisionController, condition condition.Cond, name string, handler ClusterTemplateRevisionStatusHandler) {
	statusHandler := &clusterTemplateRevisionStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromClusterTemplateRevisionHandlerToHandler(statusHandler.sync))
}

func RegisterClusterTemplateRevisionGeneratingHandler(ctx context.Context, controller ClusterTemplateRevisionController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterTemplateRevisionGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterTemplateRevisionGeneratingHandler{
		ClusterTemplateRevisionGeneratingHandler: handler,
		apply:                                    apply,
		name:                                     name,
		gvk:                                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterTemplateRevisionStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterTemplateRevisionStatusHandler struct {
	client    ClusterTemplateRevisionClient
	condition condition.Cond
	handler   ClusterTemplateRevisionStatusHandler
}

func (a *clusterTemplateRevisionStatusHandler) sync(key string, obj *v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error) {
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

type clusterTemplateRevisionGeneratingHandler struct {
	ClusterTemplateRevisionGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterTemplateRevisionGeneratingHandler) Remove(key string, obj *v3.ClusterTemplateRevision) (*v3.ClusterTemplateRevision, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.ClusterTemplateRevision{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterTemplateRevisionGeneratingHandler) Handle(obj *v3.ClusterTemplateRevision, status v3.ClusterTemplateRevisionStatus) (v3.ClusterTemplateRevisionStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ClusterTemplateRevisionGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
