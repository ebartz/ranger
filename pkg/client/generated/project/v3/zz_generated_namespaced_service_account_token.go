package client

import (
	"github.com/ranger/norman/types"
)

const (
	NamespacedServiceAccountTokenType                 = "namespacedServiceAccountToken"
	NamespacedServiceAccountTokenFieldAccountName     = "accountName"
	NamespacedServiceAccountTokenFieldAccountUID      = "accountUid"
	NamespacedServiceAccountTokenFieldAnnotations     = "annotations"
	NamespacedServiceAccountTokenFieldCACRT           = "caCrt"
	NamespacedServiceAccountTokenFieldCreated         = "created"
	NamespacedServiceAccountTokenFieldCreatorID       = "creatorId"
	NamespacedServiceAccountTokenFieldDescription     = "description"
	NamespacedServiceAccountTokenFieldLabels          = "labels"
	NamespacedServiceAccountTokenFieldName            = "name"
	NamespacedServiceAccountTokenFieldNamespaceId     = "namespaceId"
	NamespacedServiceAccountTokenFieldOwnerReferences = "ownerReferences"
	NamespacedServiceAccountTokenFieldProjectID       = "projectId"
	NamespacedServiceAccountTokenFieldRemoved         = "removed"
	NamespacedServiceAccountTokenFieldToken           = "token"
	NamespacedServiceAccountTokenFieldUUID            = "uuid"
)

type NamespacedServiceAccountToken struct {
	types.Resource
	AccountName     string            `json:"accountName,omitempty" yaml:"accountName,omitempty"`
	AccountUID      string            `json:"accountUid,omitempty" yaml:"accountUid,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CACRT           string            `json:"caCrt,omitempty" yaml:"caCrt,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId     string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ProjectID       string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	Token           string            `json:"token,omitempty" yaml:"token,omitempty"`
	UUID            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type NamespacedServiceAccountTokenCollection struct {
	types.Collection
	Data   []NamespacedServiceAccountToken `json:"data,omitempty"`
	client *NamespacedServiceAccountTokenClient
}

type NamespacedServiceAccountTokenClient struct {
	apiClient *Client
}

type NamespacedServiceAccountTokenOperations interface {
	List(opts *types.ListOpts) (*NamespacedServiceAccountTokenCollection, error)
	ListAll(opts *types.ListOpts) (*NamespacedServiceAccountTokenCollection, error)
	Create(opts *NamespacedServiceAccountToken) (*NamespacedServiceAccountToken, error)
	Update(existing *NamespacedServiceAccountToken, updates interface{}) (*NamespacedServiceAccountToken, error)
	Replace(existing *NamespacedServiceAccountToken) (*NamespacedServiceAccountToken, error)
	ByID(id string) (*NamespacedServiceAccountToken, error)
	Delete(container *NamespacedServiceAccountToken) error
}

func newNamespacedServiceAccountTokenClient(apiClient *Client) *NamespacedServiceAccountTokenClient {
	return &NamespacedServiceAccountTokenClient{
		apiClient: apiClient,
	}
}

func (c *NamespacedServiceAccountTokenClient) Create(container *NamespacedServiceAccountToken) (*NamespacedServiceAccountToken, error) {
	resp := &NamespacedServiceAccountToken{}
	err := c.apiClient.Ops.DoCreate(NamespacedServiceAccountTokenType, container, resp)
	return resp, err
}

func (c *NamespacedServiceAccountTokenClient) Update(existing *NamespacedServiceAccountToken, updates interface{}) (*NamespacedServiceAccountToken, error) {
	resp := &NamespacedServiceAccountToken{}
	err := c.apiClient.Ops.DoUpdate(NamespacedServiceAccountTokenType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *NamespacedServiceAccountTokenClient) Replace(obj *NamespacedServiceAccountToken) (*NamespacedServiceAccountToken, error) {
	resp := &NamespacedServiceAccountToken{}
	err := c.apiClient.Ops.DoReplace(NamespacedServiceAccountTokenType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *NamespacedServiceAccountTokenClient) List(opts *types.ListOpts) (*NamespacedServiceAccountTokenCollection, error) {
	resp := &NamespacedServiceAccountTokenCollection{}
	err := c.apiClient.Ops.DoList(NamespacedServiceAccountTokenType, opts, resp)
	resp.client = c
	return resp, err
}

func (c *NamespacedServiceAccountTokenClient) ListAll(opts *types.ListOpts) (*NamespacedServiceAccountTokenCollection, error) {
	resp := &NamespacedServiceAccountTokenCollection{}
	resp, err := c.List(opts)
	if err != nil {
		return resp, err
	}
	data := resp.Data
	for next, err := resp.Next(); next != nil && err == nil; next, err = next.Next() {
		data = append(data, next.Data...)
		resp = next
		resp.Data = data
	}
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (cc *NamespacedServiceAccountTokenCollection) Next() (*NamespacedServiceAccountTokenCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &NamespacedServiceAccountTokenCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *NamespacedServiceAccountTokenClient) ByID(id string) (*NamespacedServiceAccountToken, error) {
	resp := &NamespacedServiceAccountToken{}
	err := c.apiClient.Ops.DoByID(NamespacedServiceAccountTokenType, id, resp)
	return resp, err
}

func (c *NamespacedServiceAccountTokenClient) Delete(container *NamespacedServiceAccountToken) error {
	return c.apiClient.Ops.DoResourceDelete(NamespacedServiceAccountTokenType, &container.Resource)
}
