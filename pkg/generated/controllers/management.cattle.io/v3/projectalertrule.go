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

type ProjectAlertRuleHandler func(string, *v3.ProjectAlertRule) (*v3.ProjectAlertRule, error)

type ProjectAlertRuleController interface {
	generic.ControllerMeta
	ProjectAlertRuleClient

	OnChange(ctx context.Context, name string, sync ProjectAlertRuleHandler)
	OnRemove(ctx context.Context, name string, sync ProjectAlertRuleHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ProjectAlertRuleCache
}

type ProjectAlertRuleClient interface {
	Create(*v3.ProjectAlertRule) (*v3.ProjectAlertRule, error)
	Update(*v3.ProjectAlertRule) (*v3.ProjectAlertRule, error)
	UpdateStatus(*v3.ProjectAlertRule) (*v3.ProjectAlertRule, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ProjectAlertRule, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ProjectAlertRuleList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ProjectAlertRule, err error)
}

type ProjectAlertRuleCache interface {
	Get(namespace, name string) (*v3.ProjectAlertRule, error)
	List(namespace string, selector labels.Selector) ([]*v3.ProjectAlertRule, error)

	AddIndexer(indexName string, indexer ProjectAlertRuleIndexer)
	GetByIndex(indexName, key string) ([]*v3.ProjectAlertRule, error)
}

type ProjectAlertRuleIndexer func(obj *v3.ProjectAlertRule) ([]string, error)

type projectAlertRuleController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewProjectAlertRuleController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ProjectAlertRuleController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &projectAlertRuleController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromProjectAlertRuleHandlerToHandler(sync ProjectAlertRuleHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ProjectAlertRule
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ProjectAlertRule))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *projectAlertRuleController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ProjectAlertRule))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateProjectAlertRuleDeepCopyOnChange(client ProjectAlertRuleClient, obj *v3.ProjectAlertRule, handler func(obj *v3.ProjectAlertRule) (*v3.ProjectAlertRule, error)) (*v3.ProjectAlertRule, error) {
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

func (c *projectAlertRuleController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *projectAlertRuleController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *projectAlertRuleController) OnChange(ctx context.Context, name string, sync ProjectAlertRuleHandler) {
	c.AddGenericHandler(ctx, name, FromProjectAlertRuleHandlerToHandler(sync))
}

func (c *projectAlertRuleController) OnRemove(ctx context.Context, name string, sync ProjectAlertRuleHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromProjectAlertRuleHandlerToHandler(sync)))
}

func (c *projectAlertRuleController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *projectAlertRuleController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *projectAlertRuleController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *projectAlertRuleController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *projectAlertRuleController) Cache() ProjectAlertRuleCache {
	return &projectAlertRuleCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *projectAlertRuleController) Create(obj *v3.ProjectAlertRule) (*v3.ProjectAlertRule, error) {
	result := &v3.ProjectAlertRule{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *projectAlertRuleController) Update(obj *v3.ProjectAlertRule) (*v3.ProjectAlertRule, error) {
	result := &v3.ProjectAlertRule{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *projectAlertRuleController) UpdateStatus(obj *v3.ProjectAlertRule) (*v3.ProjectAlertRule, error) {
	result := &v3.ProjectAlertRule{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *projectAlertRuleController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *projectAlertRuleController) Get(namespace, name string, options metav1.GetOptions) (*v3.ProjectAlertRule, error) {
	result := &v3.ProjectAlertRule{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *projectAlertRuleController) List(namespace string, opts metav1.ListOptions) (*v3.ProjectAlertRuleList, error) {
	result := &v3.ProjectAlertRuleList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *projectAlertRuleController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *projectAlertRuleController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ProjectAlertRule, error) {
	result := &v3.ProjectAlertRule{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type projectAlertRuleCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *projectAlertRuleCache) Get(namespace, name string) (*v3.ProjectAlertRule, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ProjectAlertRule), nil
}

func (c *projectAlertRuleCache) List(namespace string, selector labels.Selector) (ret []*v3.ProjectAlertRule, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ProjectAlertRule))
	})

	return ret, err
}

func (c *projectAlertRuleCache) AddIndexer(indexName string, indexer ProjectAlertRuleIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ProjectAlertRule))
		},
	}))
}

func (c *projectAlertRuleCache) GetByIndex(indexName, key string) (result []*v3.ProjectAlertRule, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ProjectAlertRule, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ProjectAlertRule))
	}
	return result, nil
}

type ProjectAlertRuleStatusHandler func(obj *v3.ProjectAlertRule, status v3.AlertStatus) (v3.AlertStatus, error)

type ProjectAlertRuleGeneratingHandler func(obj *v3.ProjectAlertRule, status v3.AlertStatus) ([]runtime.Object, v3.AlertStatus, error)

func RegisterProjectAlertRuleStatusHandler(ctx context.Context, controller ProjectAlertRuleController, condition condition.Cond, name string, handler ProjectAlertRuleStatusHandler) {
	statusHandler := &projectAlertRuleStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromProjectAlertRuleHandlerToHandler(statusHandler.sync))
}

func RegisterProjectAlertRuleGeneratingHandler(ctx context.Context, controller ProjectAlertRuleController, apply apply.Apply,
	condition condition.Cond, name string, handler ProjectAlertRuleGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &projectAlertRuleGeneratingHandler{
		ProjectAlertRuleGeneratingHandler: handler,
		apply:                             apply,
		name:                              name,
		gvk:                               controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterProjectAlertRuleStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type projectAlertRuleStatusHandler struct {
	client    ProjectAlertRuleClient
	condition condition.Cond
	handler   ProjectAlertRuleStatusHandler
}

func (a *projectAlertRuleStatusHandler) sync(key string, obj *v3.ProjectAlertRule) (*v3.ProjectAlertRule, error) {
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

type projectAlertRuleGeneratingHandler struct {
	ProjectAlertRuleGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *projectAlertRuleGeneratingHandler) Remove(key string, obj *v3.ProjectAlertRule) (*v3.ProjectAlertRule, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.ProjectAlertRule{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *projectAlertRuleGeneratingHandler) Handle(obj *v3.ProjectAlertRule, status v3.AlertStatus) (v3.AlertStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ProjectAlertRuleGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
