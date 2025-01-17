package client

import (
	"github.com/ranger/norman/types"
)

const (
	IngressType                      = "ingress"
	IngressFieldAnnotations          = "annotations"
	IngressFieldBackend              = "backend"
	IngressFieldCreated              = "created"
	IngressFieldCreatorID            = "creatorId"
	IngressFieldDefaultBackend       = "defaultBackend"
	IngressFieldDescription          = "description"
	IngressFieldIngressClassName     = "ingressClassName"
	IngressFieldLabels               = "labels"
	IngressFieldName                 = "name"
	IngressFieldNamespaceId          = "namespaceId"
	IngressFieldOwnerReferences      = "ownerReferences"
	IngressFieldProjectID            = "projectId"
	IngressFieldPublicEndpoints      = "publicEndpoints"
	IngressFieldRemoved              = "removed"
	IngressFieldRules                = "rules"
	IngressFieldState                = "state"
	IngressFieldStatus               = "status"
	IngressFieldTLS                  = "tls"
	IngressFieldTransitioning        = "transitioning"
	IngressFieldTransitioningMessage = "transitioningMessage"
	IngressFieldUUID                 = "uuid"
)

type Ingress struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Backend              *IngressBackend   `json:"backend,omitempty" yaml:"backend,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	DefaultBackend       *IngressBackend   `json:"defaultBackend,omitempty" yaml:"defaultBackend,omitempty"`
	Description          string            `json:"description,omitempty" yaml:"description,omitempty"`
	IngressClassName     string            `json:"ingressClassName,omitempty" yaml:"ingressClassName,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId          string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ProjectID            string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	PublicEndpoints      []PublicEndpoint  `json:"publicEndpoints,omitempty" yaml:"publicEndpoints,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	Rules                []IngressRule     `json:"rules,omitempty" yaml:"rules,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Status               *IngressStatus    `json:"status,omitempty" yaml:"status,omitempty"`
	TLS                  []IngressTLS      `json:"tls,omitempty" yaml:"tls,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type IngressCollection struct {
	types.Collection
	Data   []Ingress `json:"data,omitempty"`
	client *IngressClient
}

type IngressClient struct {
	apiClient *Client
}

type IngressOperations interface {
	List(opts *types.ListOpts) (*IngressCollection, error)
	ListAll(opts *types.ListOpts) (*IngressCollection, error)
	Create(opts *Ingress) (*Ingress, error)
	Update(existing *Ingress, updates interface{}) (*Ingress, error)
	Replace(existing *Ingress) (*Ingress, error)
	ByID(id string) (*Ingress, error)
	Delete(container *Ingress) error
}

func newIngressClient(apiClient *Client) *IngressClient {
	return &IngressClient{
		apiClient: apiClient,
	}
}

func (c *IngressClient) Create(container *Ingress) (*Ingress, error) {
	resp := &Ingress{}
	err := c.apiClient.Ops.DoCreate(IngressType, container, resp)
	return resp, err
}

func (c *IngressClient) Update(existing *Ingress, updates interface{}) (*Ingress, error) {
	resp := &Ingress{}
	err := c.apiClient.Ops.DoUpdate(IngressType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *IngressClient) Replace(obj *Ingress) (*Ingress, error) {
	resp := &Ingress{}
	err := c.apiClient.Ops.DoReplace(IngressType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *IngressClient) List(opts *types.ListOpts) (*IngressCollection, error) {
	resp := &IngressCollection{}
	err := c.apiClient.Ops.DoList(IngressType, opts, resp)
	resp.client = c
	return resp, err
}

func (c *IngressClient) ListAll(opts *types.ListOpts) (*IngressCollection, error) {
	resp := &IngressCollection{}
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

func (cc *IngressCollection) Next() (*IngressCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &IngressCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *IngressClient) ByID(id string) (*Ingress, error) {
	resp := &Ingress{}
	err := c.apiClient.Ops.DoByID(IngressType, id, resp)
	return resp, err
}

func (c *IngressClient) Delete(container *Ingress) error {
	return c.apiClient.Ops.DoResourceDelete(IngressType, &container.Resource)
}
