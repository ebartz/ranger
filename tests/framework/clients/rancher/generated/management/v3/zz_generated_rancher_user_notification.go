package client

import (
	"github.com/ranger/norman/types"
)

const (
	RangerUserNotificationType                 = "rangerUserNotification"
	RangerUserNotificationFieldAnnotations     = "annotations"
	RangerUserNotificationFieldComponentName   = "componentName"
	RangerUserNotificationFieldCreated         = "created"
	RangerUserNotificationFieldCreatorID       = "creatorId"
	RangerUserNotificationFieldLabels          = "labels"
	RangerUserNotificationFieldMessage         = "message"
	RangerUserNotificationFieldName            = "name"
	RangerUserNotificationFieldOwnerReferences = "ownerReferences"
	RangerUserNotificationFieldRemoved         = "removed"
	RangerUserNotificationFieldUUID            = "uuid"
)

type RangerUserNotification struct {
	types.Resource
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ComponentName   string            `json:"componentName,omitempty" yaml:"componentName,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Message         string            `json:"message,omitempty" yaml:"message,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	UUID            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type RangerUserNotificationCollection struct {
	types.Collection
	Data   []RangerUserNotification `json:"data,omitempty"`
	client *RangerUserNotificationClient
}

type RangerUserNotificationClient struct {
	apiClient *Client
}

type RangerUserNotificationOperations interface {
	List(opts *types.ListOpts) (*RangerUserNotificationCollection, error)
	ListAll(opts *types.ListOpts) (*RangerUserNotificationCollection, error)
	Create(opts *RangerUserNotification) (*RangerUserNotification, error)
	Update(existing *RangerUserNotification, updates interface{}) (*RangerUserNotification, error)
	Replace(existing *RangerUserNotification) (*RangerUserNotification, error)
	ByID(id string) (*RangerUserNotification, error)
	Delete(container *RangerUserNotification) error
}

func newRangerUserNotificationClient(apiClient *Client) *RangerUserNotificationClient {
	return &RangerUserNotificationClient{
		apiClient: apiClient,
	}
}

func (c *RangerUserNotificationClient) Create(container *RangerUserNotification) (*RangerUserNotification, error) {
	resp := &RangerUserNotification{}
	err := c.apiClient.Ops.DoCreate(RangerUserNotificationType, container, resp)
	return resp, err
}

func (c *RangerUserNotificationClient) Update(existing *RangerUserNotification, updates interface{}) (*RangerUserNotification, error) {
	resp := &RangerUserNotification{}
	err := c.apiClient.Ops.DoUpdate(RangerUserNotificationType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *RangerUserNotificationClient) Replace(obj *RangerUserNotification) (*RangerUserNotification, error) {
	resp := &RangerUserNotification{}
	err := c.apiClient.Ops.DoReplace(RangerUserNotificationType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *RangerUserNotificationClient) List(opts *types.ListOpts) (*RangerUserNotificationCollection, error) {
	resp := &RangerUserNotificationCollection{}
	err := c.apiClient.Ops.DoList(RangerUserNotificationType, opts, resp)
	resp.client = c
	return resp, err
}

func (c *RangerUserNotificationClient) ListAll(opts *types.ListOpts) (*RangerUserNotificationCollection, error) {
	resp := &RangerUserNotificationCollection{}
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

func (cc *RangerUserNotificationCollection) Next() (*RangerUserNotificationCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &RangerUserNotificationCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *RangerUserNotificationClient) ByID(id string) (*RangerUserNotification, error) {
	resp := &RangerUserNotification{}
	err := c.apiClient.Ops.DoByID(RangerUserNotificationType, id, resp)
	return resp, err
}

func (c *RangerUserNotificationClient) Delete(container *RangerUserNotification) error {
	return c.apiClient.Ops.DoResourceDelete(RangerUserNotificationType, &container.Resource)
}
