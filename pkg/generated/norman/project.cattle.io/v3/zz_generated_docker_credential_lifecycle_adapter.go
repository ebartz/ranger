package v3

import (
	"github.com/ranger/norman/lifecycle"
	"github.com/ranger/norman/resource"
	"github.com/ranger/ranger/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type DockerCredentialLifecycle interface {
	Create(obj *v3.DockerCredential) (runtime.Object, error)
	Remove(obj *v3.DockerCredential) (runtime.Object, error)
	Updated(obj *v3.DockerCredential) (runtime.Object, error)
}

type dockerCredentialLifecycleAdapter struct {
	lifecycle DockerCredentialLifecycle
}

func (w *dockerCredentialLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *dockerCredentialLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *dockerCredentialLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v3.DockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dockerCredentialLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v3.DockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dockerCredentialLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v3.DockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewDockerCredentialLifecycleAdapter(name string, clusterScoped bool, client DockerCredentialInterface, l DockerCredentialLifecycle) DockerCredentialHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(DockerCredentialGroupVersionResource)
	}
	adapter := &dockerCredentialLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.DockerCredential) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
