package v3

import (
	"github.com/ranger/norman/lifecycle"
	"github.com/ranger/norman/resource"
	"github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type EtcdBackupLifecycle interface {
	Create(obj *v3.EtcdBackup) (runtime.Object, error)
	Remove(obj *v3.EtcdBackup) (runtime.Object, error)
	Updated(obj *v3.EtcdBackup) (runtime.Object, error)
}

type etcdBackupLifecycleAdapter struct {
	lifecycle EtcdBackupLifecycle
}

func (w *etcdBackupLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *etcdBackupLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *etcdBackupLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v3.EtcdBackup))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *etcdBackupLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v3.EtcdBackup))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *etcdBackupLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v3.EtcdBackup))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewEtcdBackupLifecycleAdapter(name string, clusterScoped bool, client EtcdBackupInterface, l EtcdBackupLifecycle) EtcdBackupHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(EtcdBackupGroupVersionResource)
	}
	adapter := &etcdBackupLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.EtcdBackup) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
