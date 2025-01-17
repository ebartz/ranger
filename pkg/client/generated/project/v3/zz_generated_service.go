package client

import (
	"github.com/ranger/norman/types"
)

const (
	ServiceType                               = "service"
	ServiceFieldAllocateLoadBalancerNodePorts = "allocateLoadBalancerNodePorts"
	ServiceFieldAnnotations                   = "annotations"
	ServiceFieldClusterIPs                    = "clusterIPs"
	ServiceFieldClusterIp                     = "clusterIp"
	ServiceFieldCreated                       = "created"
	ServiceFieldCreatorID                     = "creatorId"
	ServiceFieldDescription                   = "description"
	ServiceFieldExternalIPs                   = "externalIPs"
	ServiceFieldExternalTrafficPolicy         = "externalTrafficPolicy"
	ServiceFieldHealthCheckNodePort           = "healthCheckNodePort"
	ServiceFieldHostname                      = "hostname"
	ServiceFieldIPAddresses                   = "ipAddresses"
	ServiceFieldIPFamilies                    = "ipFamilies"
	ServiceFieldIPFamilyPolicy                = "ipFamilyPolicy"
	ServiceFieldInternalTrafficPolicy         = "internalTrafficPolicy"
	ServiceFieldKind                          = "kind"
	ServiceFieldLabels                        = "labels"
	ServiceFieldLoadBalancerClass             = "loadBalancerClass"
	ServiceFieldLoadBalancerIP                = "loadBalancerIP"
	ServiceFieldLoadBalancerSourceRanges      = "loadBalancerSourceRanges"
	ServiceFieldName                          = "name"
	ServiceFieldNamespaceId                   = "namespaceId"
	ServiceFieldOwnerReferences               = "ownerReferences"
	ServiceFieldPorts                         = "ports"
	ServiceFieldProjectID                     = "projectId"
	ServiceFieldPublicEndpoints               = "publicEndpoints"
	ServiceFieldPublishNotReadyAddresses      = "publishNotReadyAddresses"
	ServiceFieldRemoved                       = "removed"
	ServiceFieldSelector                      = "selector"
	ServiceFieldSessionAffinity               = "sessionAffinity"
	ServiceFieldSessionAffinityConfig         = "sessionAffinityConfig"
	ServiceFieldState                         = "state"
	ServiceFieldTargetDNSRecordIDs            = "targetDnsRecordIds"
	ServiceFieldTargetWorkloadIDs             = "targetWorkloadIds"
	ServiceFieldTransitioning                 = "transitioning"
	ServiceFieldTransitioningMessage          = "transitioningMessage"
	ServiceFieldUUID                          = "uuid"
	ServiceFieldWorkloadID                    = "workloadId"
)

type Service struct {
	types.Resource
	AllocateLoadBalancerNodePorts *bool                  `json:"allocateLoadBalancerNodePorts,omitempty" yaml:"allocateLoadBalancerNodePorts,omitempty"`
	Annotations                   map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ClusterIPs                    []string               `json:"clusterIPs,omitempty" yaml:"clusterIPs,omitempty"`
	ClusterIp                     string                 `json:"clusterIp,omitempty" yaml:"clusterIp,omitempty"`
	Created                       string                 `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID                     string                 `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description                   string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalIPs                   []string               `json:"externalIPs,omitempty" yaml:"externalIPs,omitempty"`
	ExternalTrafficPolicy         string                 `json:"externalTrafficPolicy,omitempty" yaml:"externalTrafficPolicy,omitempty"`
	HealthCheckNodePort           int64                  `json:"healthCheckNodePort,omitempty" yaml:"healthCheckNodePort,omitempty"`
	Hostname                      string                 `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	IPAddresses                   []string               `json:"ipAddresses,omitempty" yaml:"ipAddresses,omitempty"`
	IPFamilies                    []string               `json:"ipFamilies,omitempty" yaml:"ipFamilies,omitempty"`
	IPFamilyPolicy                string                 `json:"ipFamilyPolicy,omitempty" yaml:"ipFamilyPolicy,omitempty"`
	InternalTrafficPolicy         string                 `json:"internalTrafficPolicy,omitempty" yaml:"internalTrafficPolicy,omitempty"`
	Kind                          string                 `json:"kind,omitempty" yaml:"kind,omitempty"`
	Labels                        map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	LoadBalancerClass             string                 `json:"loadBalancerClass,omitempty" yaml:"loadBalancerClass,omitempty"`
	LoadBalancerIP                string                 `json:"loadBalancerIP,omitempty" yaml:"loadBalancerIP,omitempty"`
	LoadBalancerSourceRanges      []string               `json:"loadBalancerSourceRanges,omitempty" yaml:"loadBalancerSourceRanges,omitempty"`
	Name                          string                 `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId                   string                 `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences               []OwnerReference       `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Ports                         []ServicePort          `json:"ports,omitempty" yaml:"ports,omitempty"`
	ProjectID                     string                 `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	PublicEndpoints               []PublicEndpoint       `json:"publicEndpoints,omitempty" yaml:"publicEndpoints,omitempty"`
	PublishNotReadyAddresses      bool                   `json:"publishNotReadyAddresses,omitempty" yaml:"publishNotReadyAddresses,omitempty"`
	Removed                       string                 `json:"removed,omitempty" yaml:"removed,omitempty"`
	Selector                      map[string]string      `json:"selector,omitempty" yaml:"selector,omitempty"`
	SessionAffinity               string                 `json:"sessionAffinity,omitempty" yaml:"sessionAffinity,omitempty"`
	SessionAffinityConfig         *SessionAffinityConfig `json:"sessionAffinityConfig,omitempty" yaml:"sessionAffinityConfig,omitempty"`
	State                         string                 `json:"state,omitempty" yaml:"state,omitempty"`
	TargetDNSRecordIDs            []string               `json:"targetDnsRecordIds,omitempty" yaml:"targetDnsRecordIds,omitempty"`
	TargetWorkloadIDs             []string               `json:"targetWorkloadIds,omitempty" yaml:"targetWorkloadIds,omitempty"`
	Transitioning                 string                 `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage          string                 `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                          string                 `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	WorkloadID                    string                 `json:"workloadId,omitempty" yaml:"workloadId,omitempty"`
}

type ServiceCollection struct {
	types.Collection
	Data   []Service `json:"data,omitempty"`
	client *ServiceClient
}

type ServiceClient struct {
	apiClient *Client
}

type ServiceOperations interface {
	List(opts *types.ListOpts) (*ServiceCollection, error)
	ListAll(opts *types.ListOpts) (*ServiceCollection, error)
	Create(opts *Service) (*Service, error)
	Update(existing *Service, updates interface{}) (*Service, error)
	Replace(existing *Service) (*Service, error)
	ByID(id string) (*Service, error)
	Delete(container *Service) error
}

func newServiceClient(apiClient *Client) *ServiceClient {
	return &ServiceClient{
		apiClient: apiClient,
	}
}

func (c *ServiceClient) Create(container *Service) (*Service, error) {
	resp := &Service{}
	err := c.apiClient.Ops.DoCreate(ServiceType, container, resp)
	return resp, err
}

func (c *ServiceClient) Update(existing *Service, updates interface{}) (*Service, error) {
	resp := &Service{}
	err := c.apiClient.Ops.DoUpdate(ServiceType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ServiceClient) Replace(obj *Service) (*Service, error) {
	resp := &Service{}
	err := c.apiClient.Ops.DoReplace(ServiceType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ServiceClient) List(opts *types.ListOpts) (*ServiceCollection, error) {
	resp := &ServiceCollection{}
	err := c.apiClient.Ops.DoList(ServiceType, opts, resp)
	resp.client = c
	return resp, err
}

func (c *ServiceClient) ListAll(opts *types.ListOpts) (*ServiceCollection, error) {
	resp := &ServiceCollection{}
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

func (cc *ServiceCollection) Next() (*ServiceCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ServiceCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ServiceClient) ByID(id string) (*Service, error) {
	resp := &Service{}
	err := c.apiClient.Ops.DoByID(ServiceType, id, resp)
	return resp, err
}

func (c *ServiceClient) Delete(container *Service) error {
	return c.apiClient.Ops.DoResourceDelete(ServiceType, &container.Resource)
}
