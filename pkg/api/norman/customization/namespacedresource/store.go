package namespacedresource

import (
	"fmt"
	"strings"

	"github.com/ranger/norman/types"
	"github.com/ranger/norman/types/convert"
	"github.com/ranger/norman/types/values"
	client "github.com/ranger/ranger/pkg/client/generated/management/v3"
	v1 "github.com/ranger/ranger/pkg/generated/norman/core/v1"
	"github.com/ranger/ranger/pkg/namespace"
)

// NamespacedStore makes sure that the namespaced resources are assigned to a given namespace
type namespacedStore struct {
	types.Store
	NamespaceInterface v1.NamespaceInterface
	Namespace          string
}

func Wrap(store types.Store, nsClient v1.NamespaceInterface, namespace string) types.Store {
	return &namespacedStore{
		Store:              store,
		NamespaceInterface: nsClient,
		Namespace:          namespace,
	}
}

func (s *namespacedStore) Create(apiContext *types.APIContext, schema *types.Schema, data map[string]interface{}) (map[string]interface{}, error) {
	ns, ok := values.GetValue(data, client.PreferenceFieldNamespaceId)
	if ok && !strings.EqualFold(convert.ToString(ns), s.Namespace) {
		return nil, fmt.Errorf("error creating namespaced resource, cannot assign to %v since already assigned to %v namespace", namespace.GlobalNamespace, ns)
	} else if !ok {
		data[client.PreferenceFieldNamespaceId] = s.Namespace
	}

	return s.Store.Create(apiContext, schema, data)
}
