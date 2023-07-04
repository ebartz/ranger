package v3

import (
	"github.com/ranger/norman/lifecycle"
	"github.com/ranger/norman/resource"
	"github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type RangerUserNotificationLifecycle interface {
	Create(obj *v3.RangerUserNotification) (runtime.Object, error)
	Remove(obj *v3.RangerUserNotification) (runtime.Object, error)
	Updated(obj *v3.RangerUserNotification) (runtime.Object, error)
}

type rangerUserNotificationLifecycleAdapter struct {
	lifecycle RangerUserNotificationLifecycle
}

func (w *rangerUserNotificationLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *rangerUserNotificationLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *rangerUserNotificationLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v3.RangerUserNotification))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rangerUserNotificationLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v3.RangerUserNotification))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rangerUserNotificationLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v3.RangerUserNotification))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRangerUserNotificationLifecycleAdapter(name string, clusterScoped bool, client RangerUserNotificationInterface, l RangerUserNotificationLifecycle) RangerUserNotificationHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RangerUserNotificationGroupVersionResource)
	}
	adapter := &rangerUserNotificationLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.RangerUserNotification) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
