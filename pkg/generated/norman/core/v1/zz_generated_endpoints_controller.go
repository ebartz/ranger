package v1

import (
	"context"
	"time"

	"github.com/ranger/norman/controller"
	"github.com/ranger/norman/objectclient"
	"github.com/ranger/norman/resource"
	"k8s.io/api/core/v1"
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
	EndpointsGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Endpoints",
	}
	EndpointsResource = metav1.APIResource{
		Name:         "endpoints",
		SingularName: "endpoints",
		Namespaced:   true,

		Kind: EndpointsGroupVersionKind.Kind,
	}

	EndpointsGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "endpoints",
	}
)

func init() {
	resource.Put(EndpointsGroupVersionResource)
}

// Deprecated: use v1.Endpoints instead
type Endpoints = v1.Endpoints

func NewEndpoints(namespace, name string, obj v1.Endpoints) *v1.Endpoints {
	obj.APIVersion, obj.Kind = EndpointsGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type EndpointsHandlerFunc func(key string, obj *v1.Endpoints) (runtime.Object, error)

type EndpointsChangeHandlerFunc func(obj *v1.Endpoints) (runtime.Object, error)

type EndpointsLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1.Endpoints, err error)
	Get(namespace, name string) (*v1.Endpoints, error)
}

type EndpointsController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() EndpointsLister
	AddHandler(ctx context.Context, name string, handler EndpointsHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync EndpointsHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler EndpointsHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler EndpointsHandlerFunc)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, after time.Duration)
}

type EndpointsInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1.Endpoints) (*v1.Endpoints, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1.Endpoints, error)
	Get(name string, opts metav1.GetOptions) (*v1.Endpoints, error)
	Update(*v1.Endpoints) (*v1.Endpoints, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*v1.EndpointsList, error)
	ListNamespaced(namespace string, opts metav1.ListOptions) (*v1.EndpointsList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() EndpointsController
	AddHandler(ctx context.Context, name string, sync EndpointsHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync EndpointsHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle EndpointsLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle EndpointsLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync EndpointsHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync EndpointsHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle EndpointsLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle EndpointsLifecycle)
}

type endpointsLister struct {
	ns         string
	controller *endpointsController
}

func (l *endpointsLister) List(namespace string, selector labels.Selector) (ret []*v1.Endpoints, err error) {
	if namespace == "" {
		namespace = l.ns
	}
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1.Endpoints))
	})
	return
}

func (l *endpointsLister) Get(namespace, name string) (*v1.Endpoints, error) {
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
			Group:    EndpointsGroupVersionKind.Group,
			Resource: EndpointsGroupVersionResource.Resource,
		}, key)
	}
	return obj.(*v1.Endpoints), nil
}

type endpointsController struct {
	ns string
	controller.GenericController
}

func (c *endpointsController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *endpointsController) Lister() EndpointsLister {
	return &endpointsLister{
		ns:         c.ns,
		controller: c,
	}
}

func (c *endpointsController) AddHandler(ctx context.Context, name string, handler EndpointsHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1.Endpoints); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *endpointsController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler EndpointsHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1.Endpoints); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *endpointsController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler EndpointsHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1.Endpoints); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *endpointsController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler EndpointsHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1.Endpoints); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type endpointsFactory struct {
}

func (c endpointsFactory) Object() runtime.Object {
	return &v1.Endpoints{}
}

func (c endpointsFactory) List() runtime.Object {
	return &v1.EndpointsList{}
}

func (s *endpointsClient) Controller() EndpointsController {
	genericController := controller.NewGenericController(s.ns, EndpointsGroupVersionKind.Kind+"Controller",
		s.client.controllerFactory.ForResourceKind(EndpointsGroupVersionResource, EndpointsGroupVersionKind.Kind, true))

	return &endpointsController{
		ns:                s.ns,
		GenericController: genericController,
	}
}

type endpointsClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   EndpointsController
}

func (s *endpointsClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *endpointsClient) Create(o *v1.Endpoints) (*v1.Endpoints, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1.Endpoints), err
}

func (s *endpointsClient) Get(name string, opts metav1.GetOptions) (*v1.Endpoints, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1.Endpoints), err
}

func (s *endpointsClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1.Endpoints, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1.Endpoints), err
}

func (s *endpointsClient) Update(o *v1.Endpoints) (*v1.Endpoints, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1.Endpoints), err
}

func (s *endpointsClient) UpdateStatus(o *v1.Endpoints) (*v1.Endpoints, error) {
	obj, err := s.objectClient.UpdateStatus(o.Name, o)
	return obj.(*v1.Endpoints), err
}

func (s *endpointsClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *endpointsClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *endpointsClient) List(opts metav1.ListOptions) (*v1.EndpointsList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*v1.EndpointsList), err
}

func (s *endpointsClient) ListNamespaced(namespace string, opts metav1.ListOptions) (*v1.EndpointsList, error) {
	obj, err := s.objectClient.ListNamespaced(namespace, opts)
	return obj.(*v1.EndpointsList), err
}

func (s *endpointsClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *endpointsClient) Patch(o *v1.Endpoints, patchType types.PatchType, data []byte, subresources ...string) (*v1.Endpoints, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1.Endpoints), err
}

func (s *endpointsClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *endpointsClient) AddHandler(ctx context.Context, name string, sync EndpointsHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *endpointsClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync EndpointsHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *endpointsClient) AddLifecycle(ctx context.Context, name string, lifecycle EndpointsLifecycle) {
	sync := NewEndpointsLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *endpointsClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle EndpointsLifecycle) {
	sync := NewEndpointsLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *endpointsClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync EndpointsHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *endpointsClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync EndpointsHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *endpointsClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle EndpointsLifecycle) {
	sync := NewEndpointsLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *endpointsClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle EndpointsLifecycle) {
	sync := NewEndpointsLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}
