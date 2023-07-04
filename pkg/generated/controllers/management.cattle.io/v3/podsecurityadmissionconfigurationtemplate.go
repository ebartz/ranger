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

type PodSecurityAdmissionConfigurationTemplateHandler func(string, *v3.PodSecurityAdmissionConfigurationTemplate) (*v3.PodSecurityAdmissionConfigurationTemplate, error)

type PodSecurityAdmissionConfigurationTemplateController interface {
	generic.ControllerMeta
	PodSecurityAdmissionConfigurationTemplateClient

	OnChange(ctx context.Context, name string, sync PodSecurityAdmissionConfigurationTemplateHandler)
	OnRemove(ctx context.Context, name string, sync PodSecurityAdmissionConfigurationTemplateHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() PodSecurityAdmissionConfigurationTemplateCache
}

type PodSecurityAdmissionConfigurationTemplateClient interface {
	Create(*v3.PodSecurityAdmissionConfigurationTemplate) (*v3.PodSecurityAdmissionConfigurationTemplate, error)
	Update(*v3.PodSecurityAdmissionConfigurationTemplate) (*v3.PodSecurityAdmissionConfigurationTemplate, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.PodSecurityAdmissionConfigurationTemplate, error)
	List(opts metav1.ListOptions) (*v3.PodSecurityAdmissionConfigurationTemplateList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.PodSecurityAdmissionConfigurationTemplate, err error)
}

type PodSecurityAdmissionConfigurationTemplateCache interface {
	Get(name string) (*v3.PodSecurityAdmissionConfigurationTemplate, error)
	List(selector labels.Selector) ([]*v3.PodSecurityAdmissionConfigurationTemplate, error)

	AddIndexer(indexName string, indexer PodSecurityAdmissionConfigurationTemplateIndexer)
	GetByIndex(indexName, key string) ([]*v3.PodSecurityAdmissionConfigurationTemplate, error)
}

type PodSecurityAdmissionConfigurationTemplateIndexer func(obj *v3.PodSecurityAdmissionConfigurationTemplate) ([]string, error)

type podSecurityAdmissionConfigurationTemplateController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewPodSecurityAdmissionConfigurationTemplateController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) PodSecurityAdmissionConfigurationTemplateController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &podSecurityAdmissionConfigurationTemplateController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromPodSecurityAdmissionConfigurationTemplateHandlerToHandler(sync PodSecurityAdmissionConfigurationTemplateHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.PodSecurityAdmissionConfigurationTemplate
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.PodSecurityAdmissionConfigurationTemplate))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *podSecurityAdmissionConfigurationTemplateController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.PodSecurityAdmissionConfigurationTemplate))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdatePodSecurityAdmissionConfigurationTemplateDeepCopyOnChange(client PodSecurityAdmissionConfigurationTemplateClient, obj *v3.PodSecurityAdmissionConfigurationTemplate, handler func(obj *v3.PodSecurityAdmissionConfigurationTemplate) (*v3.PodSecurityAdmissionConfigurationTemplate, error)) (*v3.PodSecurityAdmissionConfigurationTemplate, error) {
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

func (c *podSecurityAdmissionConfigurationTemplateController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *podSecurityAdmissionConfigurationTemplateController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *podSecurityAdmissionConfigurationTemplateController) OnChange(ctx context.Context, name string, sync PodSecurityAdmissionConfigurationTemplateHandler) {
	c.AddGenericHandler(ctx, name, FromPodSecurityAdmissionConfigurationTemplateHandlerToHandler(sync))
}

func (c *podSecurityAdmissionConfigurationTemplateController) OnRemove(ctx context.Context, name string, sync PodSecurityAdmissionConfigurationTemplateHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromPodSecurityAdmissionConfigurationTemplateHandlerToHandler(sync)))
}

func (c *podSecurityAdmissionConfigurationTemplateController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *podSecurityAdmissionConfigurationTemplateController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *podSecurityAdmissionConfigurationTemplateController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *podSecurityAdmissionConfigurationTemplateController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *podSecurityAdmissionConfigurationTemplateController) Cache() PodSecurityAdmissionConfigurationTemplateCache {
	return &podSecurityAdmissionConfigurationTemplateCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *podSecurityAdmissionConfigurationTemplateController) Create(obj *v3.PodSecurityAdmissionConfigurationTemplate) (*v3.PodSecurityAdmissionConfigurationTemplate, error) {
	result := &v3.PodSecurityAdmissionConfigurationTemplate{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *podSecurityAdmissionConfigurationTemplateController) Update(obj *v3.PodSecurityAdmissionConfigurationTemplate) (*v3.PodSecurityAdmissionConfigurationTemplate, error) {
	result := &v3.PodSecurityAdmissionConfigurationTemplate{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *podSecurityAdmissionConfigurationTemplateController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *podSecurityAdmissionConfigurationTemplateController) Get(name string, options metav1.GetOptions) (*v3.PodSecurityAdmissionConfigurationTemplate, error) {
	result := &v3.PodSecurityAdmissionConfigurationTemplate{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *podSecurityAdmissionConfigurationTemplateController) List(opts metav1.ListOptions) (*v3.PodSecurityAdmissionConfigurationTemplateList, error) {
	result := &v3.PodSecurityAdmissionConfigurationTemplateList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *podSecurityAdmissionConfigurationTemplateController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *podSecurityAdmissionConfigurationTemplateController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v3.PodSecurityAdmissionConfigurationTemplate, error) {
	result := &v3.PodSecurityAdmissionConfigurationTemplate{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type podSecurityAdmissionConfigurationTemplateCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *podSecurityAdmissionConfigurationTemplateCache) Get(name string) (*v3.PodSecurityAdmissionConfigurationTemplate, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.PodSecurityAdmissionConfigurationTemplate), nil
}

func (c *podSecurityAdmissionConfigurationTemplateCache) List(selector labels.Selector) (ret []*v3.PodSecurityAdmissionConfigurationTemplate, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.PodSecurityAdmissionConfigurationTemplate))
	})

	return ret, err
}

func (c *podSecurityAdmissionConfigurationTemplateCache) AddIndexer(indexName string, indexer PodSecurityAdmissionConfigurationTemplateIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.PodSecurityAdmissionConfigurationTemplate))
		},
	}))
}

func (c *podSecurityAdmissionConfigurationTemplateCache) GetByIndex(indexName, key string) (result []*v3.PodSecurityAdmissionConfigurationTemplate, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.PodSecurityAdmissionConfigurationTemplate, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.PodSecurityAdmissionConfigurationTemplate))
	}
	return result, nil
}
