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

type ClusterAlertRuleHandler func(string, *v3.ClusterAlertRule) (*v3.ClusterAlertRule, error)

type ClusterAlertRuleController interface {
	generic.ControllerMeta
	ClusterAlertRuleClient

	OnChange(ctx context.Context, name string, sync ClusterAlertRuleHandler)
	OnRemove(ctx context.Context, name string, sync ClusterAlertRuleHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterAlertRuleCache
}

type ClusterAlertRuleClient interface {
	Create(*v3.ClusterAlertRule) (*v3.ClusterAlertRule, error)
	Update(*v3.ClusterAlertRule) (*v3.ClusterAlertRule, error)
	UpdateStatus(*v3.ClusterAlertRule) (*v3.ClusterAlertRule, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterAlertRule, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ClusterAlertRuleList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ClusterAlertRule, err error)
}

type ClusterAlertRuleCache interface {
	Get(namespace, name string) (*v3.ClusterAlertRule, error)
	List(namespace string, selector labels.Selector) ([]*v3.ClusterAlertRule, error)

	AddIndexer(indexName string, indexer ClusterAlertRuleIndexer)
	GetByIndex(indexName, key string) ([]*v3.ClusterAlertRule, error)
}

type ClusterAlertRuleIndexer func(obj *v3.ClusterAlertRule) ([]string, error)

type clusterAlertRuleController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewClusterAlertRuleController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ClusterAlertRuleController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &clusterAlertRuleController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromClusterAlertRuleHandlerToHandler(sync ClusterAlertRuleHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ClusterAlertRule
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ClusterAlertRule))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *clusterAlertRuleController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ClusterAlertRule))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateClusterAlertRuleDeepCopyOnChange(client ClusterAlertRuleClient, obj *v3.ClusterAlertRule, handler func(obj *v3.ClusterAlertRule) (*v3.ClusterAlertRule, error)) (*v3.ClusterAlertRule, error) {
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

func (c *clusterAlertRuleController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *clusterAlertRuleController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *clusterAlertRuleController) OnChange(ctx context.Context, name string, sync ClusterAlertRuleHandler) {
	c.AddGenericHandler(ctx, name, FromClusterAlertRuleHandlerToHandler(sync))
}

func (c *clusterAlertRuleController) OnRemove(ctx context.Context, name string, sync ClusterAlertRuleHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromClusterAlertRuleHandlerToHandler(sync)))
}

func (c *clusterAlertRuleController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *clusterAlertRuleController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *clusterAlertRuleController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *clusterAlertRuleController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *clusterAlertRuleController) Cache() ClusterAlertRuleCache {
	return &clusterAlertRuleCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *clusterAlertRuleController) Create(obj *v3.ClusterAlertRule) (*v3.ClusterAlertRule, error) {
	result := &v3.ClusterAlertRule{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *clusterAlertRuleController) Update(obj *v3.ClusterAlertRule) (*v3.ClusterAlertRule, error) {
	result := &v3.ClusterAlertRule{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterAlertRuleController) UpdateStatus(obj *v3.ClusterAlertRule) (*v3.ClusterAlertRule, error) {
	result := &v3.ClusterAlertRule{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterAlertRuleController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *clusterAlertRuleController) Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterAlertRule, error) {
	result := &v3.ClusterAlertRule{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *clusterAlertRuleController) List(namespace string, opts metav1.ListOptions) (*v3.ClusterAlertRuleList, error) {
	result := &v3.ClusterAlertRuleList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *clusterAlertRuleController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *clusterAlertRuleController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ClusterAlertRule, error) {
	result := &v3.ClusterAlertRule{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type clusterAlertRuleCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *clusterAlertRuleCache) Get(namespace, name string) (*v3.ClusterAlertRule, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ClusterAlertRule), nil
}

func (c *clusterAlertRuleCache) List(namespace string, selector labels.Selector) (ret []*v3.ClusterAlertRule, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ClusterAlertRule))
	})

	return ret, err
}

func (c *clusterAlertRuleCache) AddIndexer(indexName string, indexer ClusterAlertRuleIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ClusterAlertRule))
		},
	}))
}

func (c *clusterAlertRuleCache) GetByIndex(indexName, key string) (result []*v3.ClusterAlertRule, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ClusterAlertRule, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ClusterAlertRule))
	}
	return result, nil
}

type ClusterAlertRuleStatusHandler func(obj *v3.ClusterAlertRule, status v3.AlertStatus) (v3.AlertStatus, error)

type ClusterAlertRuleGeneratingHandler func(obj *v3.ClusterAlertRule, status v3.AlertStatus) ([]runtime.Object, v3.AlertStatus, error)

func RegisterClusterAlertRuleStatusHandler(ctx context.Context, controller ClusterAlertRuleController, condition condition.Cond, name string, handler ClusterAlertRuleStatusHandler) {
	statusHandler := &clusterAlertRuleStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromClusterAlertRuleHandlerToHandler(statusHandler.sync))
}

func RegisterClusterAlertRuleGeneratingHandler(ctx context.Context, controller ClusterAlertRuleController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterAlertRuleGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterAlertRuleGeneratingHandler{
		ClusterAlertRuleGeneratingHandler: handler,
		apply:                             apply,
		name:                              name,
		gvk:                               controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterAlertRuleStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterAlertRuleStatusHandler struct {
	client    ClusterAlertRuleClient
	condition condition.Cond
	handler   ClusterAlertRuleStatusHandler
}

func (a *clusterAlertRuleStatusHandler) sync(key string, obj *v3.ClusterAlertRule) (*v3.ClusterAlertRule, error) {
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

type clusterAlertRuleGeneratingHandler struct {
	ClusterAlertRuleGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterAlertRuleGeneratingHandler) Remove(key string, obj *v3.ClusterAlertRule) (*v3.ClusterAlertRule, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.ClusterAlertRule{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterAlertRuleGeneratingHandler) Handle(obj *v3.ClusterAlertRule, status v3.AlertStatus) (v3.AlertStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ClusterAlertRuleGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
