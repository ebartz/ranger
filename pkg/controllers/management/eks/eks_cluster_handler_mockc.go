//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package eks

import (
	"embed"
	"encoding/base64"
	stderrors "errors"
	"net"
	"net/url"
	"time"

	"github.com/ghodss/yaml"
	v1 "github.com/ranger/eks-operator/pkg/apis/eks.cattle.io/v1"
	"github.com/ranger/ranger/pkg/controllers/management/clusteroperator"
	mgmtv3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
	typesDialer "github.com/ranger/ranger/pkg/types/config/dialer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	secretv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

//go:embed test/*
var testFs embed.FS

type mockEksOperatorController struct {
	eksOperatorController
	mock.Mock
}

func getMockEksOperatorController(clusterState string) mockEksOperatorController {
	var dynamicClient dynamic.NamespaceableResourceInterface

	switch clusterState {
	case "default":
		dynamicClient = MockNamespaceableResourceInterfaceDefault{}
	case "create":
		dynamicClient = MockNamespaceableResourceInterfaceCreate{}
	case "active":
		dynamicClient = MockNamespaceableResourceInterfaceActive{}
	case "update":
		dynamicClient = MockNamespaceableResourceInterfaceUpdate{}
	case "Ekscc":
		dynamicClient = MockNamespaceableResourceInterfaceEksCC{}
	default:
		dynamicClient = nil
	}

	return mockEksOperatorController{
		eksOperatorController: eksOperatorController{
			OperatorController: clusteroperator.OperatorController{
				ClusterEnqueueAfter:  func(name string, duration time.Duration){},
				SecretsCache:         nil,
				Secrets:              nil,
				TemplateCache:        nil,
				ProjectCache:         nil,
				AppLister:            nil,
				AppClient:            nil,
				NsClient:             nil,
				ClusterClient:        MockClusterClient{},
				CatalogManager:       nil,
				SystemAccountManager: nil,
				DynamicClient:        dynamicClient,
				ClientDialer:         MockFactory{},
				Discovery:            MockDiscovery{},
			},
		},
		Mock: mock.Mock{},
	}
}

// test setInitialUpstreamSpec

func (m *mockEksOperatorController) setInitialUpstreamSpec(cluster *mgmtv3.Cluster) (*mgmtv3.Cluster, error) {
	logrus.Infof("setting initial upstreamSpec on cluster [%s]", cluster.Name)

	// mock
	upstreamSpec := &v1.EKSClusterConfigSpec{}

	cluster = cluster.DeepCopy()
	cluster.Status.EKSStatus.UpstreamSpec = upstreamSpec
	return m.ClusterClient.Update(cluster)
}

// test generateAndSetServiceAccount with mock sibling func (getAccessToken)

func (m *mockEksOperatorController) generateAndSetServiceAccount(cluster *mgmtv3.Cluster) (*mgmtv3.Cluster, error) {
	clusterDialer, err := m.ClientDialer.ClusterDialer(cluster.Name)
	if err != nil {
		return cluster, err
	}

	// mock
	m.Mock.On("getRestConfig", cluster).Return(&rest.Config{}, nil)

	_, err = m.getRestConfig(cluster, clusterDialer)
	if err != nil {
		return cluster, err
	}

	// mock
	secret := secretv1.Secret{}
	secret.Name = "cluster-serviceaccounttoken-sl7wm"

	cluster = cluster.DeepCopy()
	cluster.Status.ServiceAccountTokenSecret = secret.Name
	cluster.Status.ServiceAccountToken = ""
	return m.ClusterClient.Update(cluster)
}

// test generateSATokenWithPublicAPI with mock sibling func (getRestConfig)

func (m *mockEksOperatorController) generateSATokenWithPublicAPI(cluster *mgmtv3.Cluster) (string, *bool, error) {
	// mock
	m.Mock.On("getRestConfig", cluster).Return(&rest.Config{}, nil)

	_, err := m.getRestConfig(cluster, (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext)
	if err != nil {
		return "", nil, err
	}

	requiresTunnel := new(bool)

	// mock serviceToken
	serviceToken, err := "testtoken12345", nil

	if err != nil {
		*requiresTunnel = true
		var dnsError *net.DNSError
		if stderrors.As(err, &dnsError) && !dnsError.IsTemporary {
			return "", requiresTunnel, nil
		}

		// In the existence of a proxy, it may be the case that the following error occurs,
		// in which case ranger should use the tunnel connection to communicate with the cluster.
		var urlError *url.Error
		if stderrors.As(err, &urlError) && urlError.Timeout() {
			return "", requiresTunnel, nil
		}

		// Not able to determine if tunneling is required.
		requiresTunnel = nil
	}

	return serviceToken, requiresTunnel, err
}

func (m *mockEksOperatorController) getAccessToken(cluster *mgmtv3.Cluster) (string, error) {

	_mc_ret := m.Called(cluster)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func(*mgmtv3.Cluster) string); ok {
		_r0 = _rfn(cluster)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	var _r1 error

	if _rfn, ok := _mc_ret.Get(1).(func(*mgmtv3.Cluster) error); ok {
		_r1 = _rfn(cluster)
	} else {
		_r1 = _mc_ret.Error(1)
	}

	return _r0, _r1

}

func (m *mockEksOperatorController) getRestConfig(cluster *mgmtv3.Cluster, dialer typesDialer.Dialer) (*rest.Config, error) {
	// mock
	m.Mock.On("getAccessToken", cluster).Return("testaccesstoken", nil)

	accessToken, err := m.getAccessToken(cluster)
	if err != nil {
		return nil, err
	}
	decodedCA, err := base64.StdEncoding.DecodeString(cluster.Status.CACert)
	if err != nil {
		return nil, err
	}
	return &rest.Config{Host: cluster.Status.APIEndpoint, TLSClientConfig: rest.TLSClientConfig{CAData: decodedCA}, BearerToken: accessToken, Dial: dialer}, nil
}

// utility

func getMockV3Cluster(filename string) (mgmtv3.Cluster, error) {
	var mockCluster mgmtv3.Cluster

	// Read the embedded file
	cluster, err := testFs.ReadFile(filename); if err != nil {
		return mockCluster, err
	}
	// Unmarshal cluster yaml into a management v3 cluster object
	err = yaml.Unmarshal(cluster, &mockCluster); if err != nil {
		return mockCluster, err
	}

	return mockCluster, nil
}

func getMockEksClusterConfig(filename string) (*unstructured.Unstructured, error) {
	var EksClusterConfig *unstructured.Unstructured

	// Read the embedded file
	bytes, err := testFs.ReadFile(filename); if err != nil {
		return EksClusterConfig, err
	}
	// Unmarshal json into an unstructured cluster config object
	err = json.Unmarshal(bytes, &EksClusterConfig); if err != nil {
		return EksClusterConfig, err
	}

	return EksClusterConfig, nil
}
