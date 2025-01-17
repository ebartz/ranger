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
	PodSecurityPolicyTemplateProjectBindingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "PodSecurityPolicyTemplateProjectBinding",
	}
	PodSecurityPolicyTemplateProjectBindingResource = metav1.APIResource{
		Name:         "podsecuritypolicytemplateprojectbindings",
		SingularName: "podsecuritypolicytemplateprojectbinding",
		Namespaced:   true,

		Kind: PodSecurityPolicyTemplateProjectBindingGroupVersionKind.Kind,
	}

	PodSecurityPolicyTemplateProjectBindingGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "podsecuritypolicytemplateprojectbindings",
	}
)

func init() {
	resource.Put(PodSecurityPolicyTemplateProjectBindingGroupVersionResource)
}

// Deprecated: use v3.PodSecurityPolicyTemplateProjectBinding instead
type PodSecurityPolicyTemplateProjectBinding = v3.PodSecurityPolicyTemplateProjectBinding

func NewPodSecurityPolicyTemplateProjectBinding(namespace, name string, obj v3.PodSecurityPolicyTemplateProjectBinding) *v3.PodSecurityPolicyTemplateProjectBinding {
	obj.APIVersion, obj.Kind = PodSecurityPolicyTemplateProjectBindingGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type PodSecurityPolicyTemplateProjectBindingHandlerFunc func(key string, obj *v3.PodSecurityPolicyTemplateProjectBinding) (runtime.Object, error)

type PodSecurityPolicyTemplateProjectBindingChangeHandlerFunc func(obj *v3.PodSecurityPolicyTemplateProjectBinding) (runtime.Object, error)

type PodSecurityPolicyTemplateProjectBindingLister interface {
	List(namespace string, selector labels.Selector) (ret []*v3.PodSecurityPolicyTemplateProjectBinding, err error)
	Get(namespace, name string) (*v3.PodSecurityPolicyTemplateProjectBinding, error)
}

type PodSecurityPolicyTemplateProjectBindingController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() PodSecurityPolicyTemplateProjectBindingLister
	AddHandler(ctx context.Context, name string, handler PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, after time.Duration)
}

type PodSecurityPolicyTemplateProjectBindingInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v3.PodSecurityPolicyTemplateProjectBinding) (*v3.PodSecurityPolicyTemplateProjectBinding, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.PodSecurityPolicyTemplateProjectBinding, error)
	Get(name string, opts metav1.GetOptions) (*v3.PodSecurityPolicyTemplateProjectBinding, error)
	Update(*v3.PodSecurityPolicyTemplateProjectBinding) (*v3.PodSecurityPolicyTemplateProjectBinding, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*v3.PodSecurityPolicyTemplateProjectBindingList, error)
	ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.PodSecurityPolicyTemplateProjectBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() PodSecurityPolicyTemplateProjectBindingController
	AddHandler(ctx context.Context, name string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle)
}

type podSecurityPolicyTemplateProjectBindingLister struct {
	ns         string
	controller *podSecurityPolicyTemplateProjectBindingController
}

func (l *podSecurityPolicyTemplateProjectBindingLister) List(namespace string, selector labels.Selector) (ret []*v3.PodSecurityPolicyTemplateProjectBinding, err error) {
	if namespace == "" {
		namespace = l.ns
	}
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v3.PodSecurityPolicyTemplateProjectBinding))
	})
	return
}

func (l *podSecurityPolicyTemplateProjectBindingLister) Get(namespace, name string) (*v3.PodSecurityPolicyTemplateProjectBinding, error) {
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
			Group:    PodSecurityPolicyTemplateProjectBindingGroupVersionKind.Group,
			Resource: PodSecurityPolicyTemplateProjectBindingGroupVersionResource.Resource,
		}, key)
	}
	return obj.(*v3.PodSecurityPolicyTemplateProjectBinding), nil
}

type podSecurityPolicyTemplateProjectBindingController struct {
	ns string
	controller.GenericController
}

func (c *podSecurityPolicyTemplateProjectBindingController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *podSecurityPolicyTemplateProjectBindingController) Lister() PodSecurityPolicyTemplateProjectBindingLister {
	return &podSecurityPolicyTemplateProjectBindingLister{
		ns:         c.ns,
		controller: c,
	}
}

func (c *podSecurityPolicyTemplateProjectBindingController) AddHandler(ctx context.Context, name string, handler PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.PodSecurityPolicyTemplateProjectBinding); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *podSecurityPolicyTemplateProjectBindingController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.PodSecurityPolicyTemplateProjectBinding); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *podSecurityPolicyTemplateProjectBindingController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.PodSecurityPolicyTemplateProjectBinding); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *podSecurityPolicyTemplateProjectBindingController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.PodSecurityPolicyTemplateProjectBinding); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type podSecurityPolicyTemplateProjectBindingFactory struct {
}

func (c podSecurityPolicyTemplateProjectBindingFactory) Object() runtime.Object {
	return &v3.PodSecurityPolicyTemplateProjectBinding{}
}

func (c podSecurityPolicyTemplateProjectBindingFactory) List() runtime.Object {
	return &v3.PodSecurityPolicyTemplateProjectBindingList{}
}

func (s *podSecurityPolicyTemplateProjectBindingClient) Controller() PodSecurityPolicyTemplateProjectBindingController {
	genericController := controller.NewGenericController(s.ns, PodSecurityPolicyTemplateProjectBindingGroupVersionKind.Kind+"Controller",
		s.client.controllerFactory.ForResourceKind(PodSecurityPolicyTemplateProjectBindingGroupVersionResource, PodSecurityPolicyTemplateProjectBindingGroupVersionKind.Kind, true))

	return &podSecurityPolicyTemplateProjectBindingController{
		ns:                s.ns,
		GenericController: genericController,
	}
}

type podSecurityPolicyTemplateProjectBindingClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   PodSecurityPolicyTemplateProjectBindingController
}

func (s *podSecurityPolicyTemplateProjectBindingClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *podSecurityPolicyTemplateProjectBindingClient) Create(o *v3.PodSecurityPolicyTemplateProjectBinding) (*v3.PodSecurityPolicyTemplateProjectBinding, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBinding), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) Get(name string, opts metav1.GetOptions) (*v3.PodSecurityPolicyTemplateProjectBinding, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBinding), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.PodSecurityPolicyTemplateProjectBinding, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBinding), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) Update(o *v3.PodSecurityPolicyTemplateProjectBinding) (*v3.PodSecurityPolicyTemplateProjectBinding, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBinding), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) UpdateStatus(o *v3.PodSecurityPolicyTemplateProjectBinding) (*v3.PodSecurityPolicyTemplateProjectBinding, error) {
	obj, err := s.objectClient.UpdateStatus(o.Name, o)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBinding), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) List(opts metav1.ListOptions) (*v3.PodSecurityPolicyTemplateProjectBindingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBindingList), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.PodSecurityPolicyTemplateProjectBindingList, error) {
	obj, err := s.objectClient.ListNamespaced(namespace, opts)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBindingList), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *podSecurityPolicyTemplateProjectBindingClient) Patch(o *v3.PodSecurityPolicyTemplateProjectBinding, patchType types.PatchType, data []byte, subresources ...string) (*v3.PodSecurityPolicyTemplateProjectBinding, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v3.PodSecurityPolicyTemplateProjectBinding), err
}

func (s *podSecurityPolicyTemplateProjectBindingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddHandler(ctx context.Context, name string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddLifecycle(ctx context.Context, name string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle) {
	sync := NewPodSecurityPolicyTemplateProjectBindingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle) {
	sync := NewPodSecurityPolicyTemplateProjectBindingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync PodSecurityPolicyTemplateProjectBindingHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle) {
	sync := NewPodSecurityPolicyTemplateProjectBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *podSecurityPolicyTemplateProjectBindingClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle PodSecurityPolicyTemplateProjectBindingLifecycle) {
	sync := NewPodSecurityPolicyTemplateProjectBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}
