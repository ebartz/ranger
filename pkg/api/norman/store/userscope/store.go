package userscope

import (
	"fmt"
	"strings"

	"github.com/ranger/norman/httperror"
	"github.com/ranger/norman/types"
	client "github.com/ranger/ranger/pkg/client/generated/management/v3"
	v1 "github.com/ranger/ranger/pkg/generated/norman/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NamespaceID = client.PreferenceFieldNamespaceId
)

type Store struct {
	Store    types.Store
	nsClient v1.NamespaceInterface
}

func NewStore(nsClient v1.NamespaceInterface, store types.Store) *Store {
	return &Store{
		Store:    store,
		nsClient: nsClient,
	}
}

func (s *Store) Context() types.StorageContext {
	return s.Store.Context()
}

func (s *Store) ByID(apiContext *types.APIContext, schema *types.Schema, id string) (map[string]interface{}, error) {
	user, err := getUser(apiContext)
	if err != nil {
		return nil, err
	}

	return s.Store.ByID(apiContext, schema, addNamespace(user, id))
}

func (s *Store) List(apiContext *types.APIContext, schema *types.Schema, opt *types.QueryOptions) ([]map[string]interface{}, error) {
	user, err := getUser(apiContext)
	if err != nil {
		return nil, err
	}

	if opt == nil {
		return nil, nil
	}

	opt.Conditions = append(opt.Conditions, types.EQ(NamespaceID, getNamespace(user)))
	return s.Store.List(apiContext, schema, opt)
}

func (s *Store) Create(apiContext *types.APIContext, schema *types.Schema, data map[string]interface{}) (map[string]interface{}, error) {
	user, err := getUser(apiContext)
	if err != nil || data == nil {
		return nil, err
	}

	ns := getNamespace(user)
	_, err = s.nsClient.Get(ns, metav1.GetOptions{})
	if err != nil {
		s.nsClient.Create(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"management.cattle.io/system-namespace": "true",
				},
				Name: ns,
			},
		})
	}

	data[NamespaceID] = getNamespace(user)
	return s.Store.Create(apiContext, schema, data)
}

func (s *Store) Update(apiContext *types.APIContext, schema *types.Schema, data map[string]interface{}, id string) (map[string]interface{}, error) {
	user, err := getUser(apiContext)
	if err != nil {
		return nil, err
	}

	return s.Store.Update(apiContext, schema, data, addNamespace(user, id))
}

func (s *Store) Delete(apiContext *types.APIContext, schema *types.Schema, id string) (map[string]interface{}, error) {
	user, err := getUser(apiContext)
	if err != nil {
		return nil, err
	}

	return s.Store.Delete(apiContext, schema, addNamespace(user, id))
}

func (s *Store) Watch(apiContext *types.APIContext, schema *types.Schema, opt *types.QueryOptions) (chan map[string]interface{}, error) {
	user, err := getUser(apiContext)
	if err != nil {
		return nil, err
	}

	if opt == nil {
		return nil, nil
	}

	opt.Conditions = append(opt.Conditions, types.EQ(NamespaceID, getNamespace(user)))
	return s.Store.Watch(apiContext, schema, opt)
}

func getUser(apiContext *types.APIContext) (string, error) {
	user := apiContext.Request.Header.Get("Impersonate-User")
	if user == "" {
		return "", httperror.NewAPIError(httperror.NotFound, "missing user")
	}
	return user, nil
}

func addNamespace(user, id string) string {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 1 {
		return fmt.Sprintf("%s:%s", getNamespace(user), parts[0])
	}
	return fmt.Sprintf("%s:%s", getNamespace(user), parts[1])
}

func getNamespace(user string) string {
	return user
}
