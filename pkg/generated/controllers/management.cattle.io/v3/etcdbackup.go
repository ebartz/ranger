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

type EtcdBackupHandler func(string, *v3.EtcdBackup) (*v3.EtcdBackup, error)

type EtcdBackupController interface {
	generic.ControllerMeta
	EtcdBackupClient

	OnChange(ctx context.Context, name string, sync EtcdBackupHandler)
	OnRemove(ctx context.Context, name string, sync EtcdBackupHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() EtcdBackupCache
}

type EtcdBackupClient interface {
	Create(*v3.EtcdBackup) (*v3.EtcdBackup, error)
	Update(*v3.EtcdBackup) (*v3.EtcdBackup, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.EtcdBackup, error)
	List(namespace string, opts metav1.ListOptions) (*v3.EtcdBackupList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.EtcdBackup, err error)
}

type EtcdBackupCache interface {
	Get(namespace, name string) (*v3.EtcdBackup, error)
	List(namespace string, selector labels.Selector) ([]*v3.EtcdBackup, error)

	AddIndexer(indexName string, indexer EtcdBackupIndexer)
	GetByIndex(indexName, key string) ([]*v3.EtcdBackup, error)
}

type EtcdBackupIndexer func(obj *v3.EtcdBackup) ([]string, error)

type etcdBackupController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewEtcdBackupController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) EtcdBackupController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &etcdBackupController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromEtcdBackupHandlerToHandler(sync EtcdBackupHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.EtcdBackup
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.EtcdBackup))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *etcdBackupController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.EtcdBackup))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateEtcdBackupDeepCopyOnChange(client EtcdBackupClient, obj *v3.EtcdBackup, handler func(obj *v3.EtcdBackup) (*v3.EtcdBackup, error)) (*v3.EtcdBackup, error) {
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

func (c *etcdBackupController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *etcdBackupController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *etcdBackupController) OnChange(ctx context.Context, name string, sync EtcdBackupHandler) {
	c.AddGenericHandler(ctx, name, FromEtcdBackupHandlerToHandler(sync))
}

func (c *etcdBackupController) OnRemove(ctx context.Context, name string, sync EtcdBackupHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromEtcdBackupHandlerToHandler(sync)))
}

func (c *etcdBackupController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *etcdBackupController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *etcdBackupController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *etcdBackupController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *etcdBackupController) Cache() EtcdBackupCache {
	return &etcdBackupCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *etcdBackupController) Create(obj *v3.EtcdBackup) (*v3.EtcdBackup, error) {
	result := &v3.EtcdBackup{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *etcdBackupController) Update(obj *v3.EtcdBackup) (*v3.EtcdBackup, error) {
	result := &v3.EtcdBackup{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *etcdBackupController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *etcdBackupController) Get(namespace, name string, options metav1.GetOptions) (*v3.EtcdBackup, error) {
	result := &v3.EtcdBackup{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *etcdBackupController) List(namespace string, opts metav1.ListOptions) (*v3.EtcdBackupList, error) {
	result := &v3.EtcdBackupList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *etcdBackupController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *etcdBackupController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.EtcdBackup, error) {
	result := &v3.EtcdBackup{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type etcdBackupCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *etcdBackupCache) Get(namespace, name string) (*v3.EtcdBackup, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.EtcdBackup), nil
}

func (c *etcdBackupCache) List(namespace string, selector labels.Selector) (ret []*v3.EtcdBackup, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.EtcdBackup))
	})

	return ret, err
}

func (c *etcdBackupCache) AddIndexer(indexName string, indexer EtcdBackupIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.EtcdBackup))
		},
	}))
}

func (c *etcdBackupCache) GetByIndex(indexName, key string) (result []*v3.EtcdBackup, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.EtcdBackup, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.EtcdBackup))
	}
	return result, nil
}
