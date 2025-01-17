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
	ProjectNetworkPolicyGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ProjectNetworkPolicy",
	}
	ProjectNetworkPolicyResource = metav1.APIResource{
		Name:         "projectnetworkpolicies",
		SingularName: "projectnetworkpolicy",
		Namespaced:   true,

		Kind: ProjectNetworkPolicyGroupVersionKind.Kind,
	}

	ProjectNetworkPolicyGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "projectnetworkpolicies",
	}
)

func init() {
	resource.Put(ProjectNetworkPolicyGroupVersionResource)
}

// Deprecated: use v3.ProjectNetworkPolicy instead
type ProjectNetworkPolicy = v3.ProjectNetworkPolicy

func NewProjectNetworkPolicy(namespace, name string, obj v3.ProjectNetworkPolicy) *v3.ProjectNetworkPolicy {
	obj.APIVersion, obj.Kind = ProjectNetworkPolicyGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ProjectNetworkPolicyHandlerFunc func(key string, obj *v3.ProjectNetworkPolicy) (runtime.Object, error)

type ProjectNetworkPolicyChangeHandlerFunc func(obj *v3.ProjectNetworkPolicy) (runtime.Object, error)

type ProjectNetworkPolicyLister interface {
	List(namespace string, selector labels.Selector) (ret []*v3.ProjectNetworkPolicy, err error)
	Get(namespace, name string) (*v3.ProjectNetworkPolicy, error)
}

type ProjectNetworkPolicyController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ProjectNetworkPolicyLister
	AddHandler(ctx context.Context, name string, handler ProjectNetworkPolicyHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ProjectNetworkPolicyHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ProjectNetworkPolicyHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler ProjectNetworkPolicyHandlerFunc)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, after time.Duration)
}

type ProjectNetworkPolicyInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.ProjectNetworkPolicy, error)
	Get(name string, opts metav1.GetOptions) (*v3.ProjectNetworkPolicy, error)
	Update(*v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*v3.ProjectNetworkPolicyList, error)
	ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.ProjectNetworkPolicyList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ProjectNetworkPolicyController
	AddHandler(ctx context.Context, name string, sync ProjectNetworkPolicyHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ProjectNetworkPolicyHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ProjectNetworkPolicyLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ProjectNetworkPolicyLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ProjectNetworkPolicyHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ProjectNetworkPolicyHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ProjectNetworkPolicyLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ProjectNetworkPolicyLifecycle)
}

type projectNetworkPolicyLister struct {
	ns         string
	controller *projectNetworkPolicyController
}

func (l *projectNetworkPolicyLister) List(namespace string, selector labels.Selector) (ret []*v3.ProjectNetworkPolicy, err error) {
	if namespace == "" {
		namespace = l.ns
	}
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v3.ProjectNetworkPolicy))
	})
	return
}

func (l *projectNetworkPolicyLister) Get(namespace, name string) (*v3.ProjectNetworkPolicy, error) {
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
			Group:    ProjectNetworkPolicyGroupVersionKind.Group,
			Resource: ProjectNetworkPolicyGroupVersionResource.Resource,
		}, key)
	}
	return obj.(*v3.ProjectNetworkPolicy), nil
}

type projectNetworkPolicyController struct {
	ns string
	controller.GenericController
}

func (c *projectNetworkPolicyController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *projectNetworkPolicyController) Lister() ProjectNetworkPolicyLister {
	return &projectNetworkPolicyLister{
		ns:         c.ns,
		controller: c,
	}
}

func (c *projectNetworkPolicyController) AddHandler(ctx context.Context, name string, handler ProjectNetworkPolicyHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.ProjectNetworkPolicy); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *projectNetworkPolicyController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler ProjectNetworkPolicyHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.ProjectNetworkPolicy); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *projectNetworkPolicyController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ProjectNetworkPolicyHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.ProjectNetworkPolicy); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *projectNetworkPolicyController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler ProjectNetworkPolicyHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.ProjectNetworkPolicy); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type projectNetworkPolicyFactory struct {
}

func (c projectNetworkPolicyFactory) Object() runtime.Object {
	return &v3.ProjectNetworkPolicy{}
}

func (c projectNetworkPolicyFactory) List() runtime.Object {
	return &v3.ProjectNetworkPolicyList{}
}

func (s *projectNetworkPolicyClient) Controller() ProjectNetworkPolicyController {
	genericController := controller.NewGenericController(s.ns, ProjectNetworkPolicyGroupVersionKind.Kind+"Controller",
		s.client.controllerFactory.ForResourceKind(ProjectNetworkPolicyGroupVersionResource, ProjectNetworkPolicyGroupVersionKind.Kind, true))

	return &projectNetworkPolicyController{
		ns:                s.ns,
		GenericController: genericController,
	}
}

type projectNetworkPolicyClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ProjectNetworkPolicyController
}

func (s *projectNetworkPolicyClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *projectNetworkPolicyClient) Create(o *v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v3.ProjectNetworkPolicy), err
}

func (s *projectNetworkPolicyClient) Get(name string, opts metav1.GetOptions) (*v3.ProjectNetworkPolicy, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v3.ProjectNetworkPolicy), err
}

func (s *projectNetworkPolicyClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.ProjectNetworkPolicy, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v3.ProjectNetworkPolicy), err
}

func (s *projectNetworkPolicyClient) Update(o *v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v3.ProjectNetworkPolicy), err
}

func (s *projectNetworkPolicyClient) UpdateStatus(o *v3.ProjectNetworkPolicy) (*v3.ProjectNetworkPolicy, error) {
	obj, err := s.objectClient.UpdateStatus(o.Name, o)
	return obj.(*v3.ProjectNetworkPolicy), err
}

func (s *projectNetworkPolicyClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *projectNetworkPolicyClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *projectNetworkPolicyClient) List(opts metav1.ListOptions) (*v3.ProjectNetworkPolicyList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*v3.ProjectNetworkPolicyList), err
}

func (s *projectNetworkPolicyClient) ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.ProjectNetworkPolicyList, error) {
	obj, err := s.objectClient.ListNamespaced(namespace, opts)
	return obj.(*v3.ProjectNetworkPolicyList), err
}

func (s *projectNetworkPolicyClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *projectNetworkPolicyClient) Patch(o *v3.ProjectNetworkPolicy, patchType types.PatchType, data []byte, subresources ...string) (*v3.ProjectNetworkPolicy, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v3.ProjectNetworkPolicy), err
}

func (s *projectNetworkPolicyClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *projectNetworkPolicyClient) AddHandler(ctx context.Context, name string, sync ProjectNetworkPolicyHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *projectNetworkPolicyClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync ProjectNetworkPolicyHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *projectNetworkPolicyClient) AddLifecycle(ctx context.Context, name string, lifecycle ProjectNetworkPolicyLifecycle) {
	sync := NewProjectNetworkPolicyLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *projectNetworkPolicyClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle ProjectNetworkPolicyLifecycle) {
	sync := NewProjectNetworkPolicyLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *projectNetworkPolicyClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ProjectNetworkPolicyHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *projectNetworkPolicyClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync ProjectNetworkPolicyHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *projectNetworkPolicyClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ProjectNetworkPolicyLifecycle) {
	sync := NewProjectNetworkPolicyLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *projectNetworkPolicyClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle ProjectNetworkPolicyLifecycle) {
	sync := NewProjectNetworkPolicyLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}
