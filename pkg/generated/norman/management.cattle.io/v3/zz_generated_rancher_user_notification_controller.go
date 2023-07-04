package v3

import (
	"context"
	"time"

	"github.com/ranger/norman/controller"
	"github.com/ranger/norman/objectclient"
	"github.com/ranger/norman/resource"
	"github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	RangerUserNotificationGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "RangerUserNotification",
	}
	RangerUserNotificationResource = metav1.APIResource{
		Name:         "rangerusernotifications",
		SingularName: "rangerusernotification",
		Namespaced:   false,
		Kind:         RangerUserNotificationGroupVersionKind.Kind,
	}

	RangerUserNotificationGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "rangerusernotifications",
	}
)

func init() {
	resource.Put(RangerUserNotificationGroupVersionResource)
}

// Deprecated: use v3.RangerUserNotification instead
type RangerUserNotification = v3.RangerUserNotification

func NewRangerUserNotification(namespace, name string, obj v3.RangerUserNotification) *v3.RangerUserNotification {
	obj.APIVersion, obj.Kind = RangerUserNotificationGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type RangerUserNotificationHandlerFunc func(key string, obj *v3.RangerUserNotification) (runtime.Object, error)

type RangerUserNotificationChangeHandlerFunc func(obj *v3.RangerUserNotification) (runtime.Object, error)

type RangerUserNotificationLister interface {
	List(namespace string, selector labels.Selector) (ret []*v3.RangerUserNotification, err error)
	Get(namespace, name string) (*v3.RangerUserNotification, error)
}

type RangerUserNotificationController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() RangerUserNotificationLister
	AddHandler(ctx context.Context, name string, handler RangerUserNotificationHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync RangerUserNotificationHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler RangerUserNotificationHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler RangerUserNotificationHandlerFunc)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, after time.Duration)
}

type RangerUserNotificationInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v3.RangerUserNotification) (*v3.RangerUserNotification, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.RangerUserNotification, error)
	Get(name string, opts metav1.GetOptions) (*v3.RangerUserNotification, error)
	Update(*v3.RangerUserNotification) (*v3.RangerUserNotification, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*v3.RangerUserNotificationList, error)
	ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.RangerUserNotificationList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() RangerUserNotificationController
	AddHandler(ctx context.Context, name string, sync RangerUserNotificationHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync RangerUserNotificationHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle RangerUserNotificationLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle RangerUserNotificationLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync RangerUserNotificationHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync RangerUserNotificationHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle RangerUserNotificationLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle RangerUserNotificationLifecycle)
}

type rangerUserNotificationLister struct {
	ns         string
	controller *rangerUserNotificationController
}

func (l *rangerUserNotificationLister) List(namespace string, selector labels.Selector) (ret []*v3.RangerUserNotification, err error) {
	if namespace == "" {
		namespace = l.ns
	}
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v3.RangerUserNotification))
	})
	return
}

func (l *rangerUserNotificationLister) Get(namespace, name string) (*v3.RangerUserNotification, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    RangerUserNotificationGroupVersionKind.Group,
			Resource: RangerUserNotificationGroupVersionResource.Resource,
		}, key)
	}
	return obj.(*v3.RangerUserNotification), nil
}

type rangerUserNotificationController struct {
	ns string
	controller.GenericController
}

func (c *rangerUserNotificationController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *rangerUserNotificationController) Lister() RangerUserNotificationLister {
	return &rangerUserNotificationLister{
		ns:         c.ns,
		controller: c,
	}
}

func (c *rangerUserNotificationController) AddHandler(ctx context.Context, name string, handler RangerUserNotificationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.RangerUserNotification); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *rangerUserNotificationController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler RangerUserNotificationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.RangerUserNotification); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *rangerUserNotificationController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler RangerUserNotificationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.RangerUserNotification); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *rangerUserNotificationController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler RangerUserNotificationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.RangerUserNotification); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type rangerUserNotificationFactory struct {
}

func (c rangerUserNotificationFactory) Object() runtime.Object {
	return &v3.RangerUserNotification{}
}

func (c rangerUserNotificationFactory) List() runtime.Object {
	return &v3.RangerUserNotificationList{}
}

func (s *rangerUserNotificationClient) Controller() RangerUserNotificationController {
	genericController := controller.NewGenericController(s.ns, RangerUserNotificationGroupVersionKind.Kind+"Controller",
		s.client.controllerFactory.ForResourceKind(RangerUserNotificationGroupVersionResource, RangerUserNotificationGroupVersionKind.Kind, false))

	return &rangerUserNotificationController{
		ns:                s.ns,
		GenericController: genericController,
	}
}

type rangerUserNotificationClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   RangerUserNotificationController
}

func (s *rangerUserNotificationClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *rangerUserNotificationClient) Create(o *v3.RangerUserNotification) (*v3.RangerUserNotification, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v3.RangerUserNotification), err
}

func (s *rangerUserNotificationClient) Get(name string, opts metav1.GetOptions) (*v3.RangerUserNotification, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v3.RangerUserNotification), err
}

func (s *rangerUserNotificationClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.RangerUserNotification, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v3.RangerUserNotification), err
}

func (s *rangerUserNotificationClient) Update(o *v3.RangerUserNotification) (*v3.RangerUserNotification, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v3.RangerUserNotification), err
}

func (s *rangerUserNotificationClient) UpdateStatus(o *v3.RangerUserNotification) (*v3.RangerUserNotification, error) {
	obj, err := s.objectClient.UpdateStatus(o.Name, o)
	return obj.(*v3.RangerUserNotification), err
}

func (s *rangerUserNotificationClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *rangerUserNotificationClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *rangerUserNotificationClient) List(opts metav1.ListOptions) (*v3.RangerUserNotificationList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*v3.RangerUserNotificationList), err
}

func (s *rangerUserNotificationClient) ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.RangerUserNotificationList, error) {
	obj, err := s.objectClient.ListNamespaced(namespace, opts)
	return obj.(*v3.RangerUserNotificationList), err
}

func (s *rangerUserNotificationClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *rangerUserNotificationClient) Patch(o *v3.RangerUserNotification, patchType types.PatchType, data []byte, subresources ...string) (*v3.RangerUserNotification, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v3.RangerUserNotification), err
}

func (s *rangerUserNotificationClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *rangerUserNotificationClient) AddHandler(ctx context.Context, name string, sync RangerUserNotificationHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *rangerUserNotificationClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync RangerUserNotificationHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *rangerUserNotificationClient) AddLifecycle(ctx context.Context, name string, lifecycle RangerUserNotificationLifecycle) {
	sync := NewRangerUserNotificationLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *rangerUserNotificationClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle RangerUserNotificationLifecycle) {
	sync := NewRangerUserNotificationLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *rangerUserNotificationClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync RangerUserNotificationHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *rangerUserNotificationClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync RangerUserNotificationHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *rangerUserNotificationClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle RangerUserNotificationLifecycle) {
	sync := NewRangerUserNotificationLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *rangerUserNotificationClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle RangerUserNotificationLifecycle) {
	sync := NewRangerUserNotificationLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}
