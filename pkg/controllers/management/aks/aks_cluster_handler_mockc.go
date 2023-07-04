//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package aks

import (
	"embed"
	stderrors "errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/ghodss/yaml"
	v1 "github.com/ranger/aks-operator/pkg/apis/aks.cattle.io/v1"
	"github.com/ranger/ranger/pkg/controllers/management/clusteroperator"
	mgmtv3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
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

type mockAksOperatorController struct {
	aksOperatorController
	mock.Mock
}

func getMockAksOperatorController(clusterState string) mockAksOperatorController {
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
	case "akscc":
		dynamicClient = MockNamespaceableResourceInterfaceAKSCC{}
	default:
		dynamicClient = nil
	}

	return mockAksOperatorController{
		aksOperatorController: aksOperatorController{
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
				ClientDialer: 		  MockFactory{},
				Discovery:            MockDiscovery{},
			},
			secretClient: nil,
		},
		Mock: mock.Mock{},
	}
}

// test setInitialUpstreamSpec

func (m *mockAksOperatorController) setInitialUpstreamSpec(cluster *mgmtv3.Cluster) (*mgmtv3.Cluster, error) {
	logrus.Infof("setting initial upstreamSpec on cluster [%s]", cluster.Name)

	// mock
	upstreamSpec := &v1.AKSClusterConfigSpec{}

	cluster = cluster.DeepCopy()
	cluster.Status.AKSStatus.UpstreamSpec = upstreamSpec
	return m.ClusterClient.Update(cluster)
}

// test generateAndSetServiceAccount with mock sibling func (getRestConfig)

func (m *mockAksOperatorController) generateAndSetServiceAccount(cluster *mgmtv3.Cluster) (*mgmtv3.Cluster, error) {
	// mock
	m.Mock.On("getRestConfig", cluster).Return(&rest.Config{}, nil)

	restConfig, err := m.getRestConfig(cluster)
	if err != nil {
		return cluster, fmt.Errorf("error getting kube config: %v", err)
	}

	clusterDialer, err := m.ClientDialer.ClusterDialer(cluster.Name)
	if err != nil {
		return cluster, err
	}

	restConfig.Dial = clusterDialer
	cluster = cluster.DeepCopy()

	// mock
	secret := secretv1.Secret{}
	secret.Name = "cluster-serviceaccounttoken-sl7wm"

	if err != nil {
		return nil, err
	}
	cluster.Status.ServiceAccountTokenSecret = secret.Name
	cluster.Status.ServiceAccountToken = ""
	return m.ClusterClient.Update(cluster)
}

// test generateSATokenWithPublicAPI with mock sibling func (getRestConfig)

func (m *mockAksOperatorController) generateSATokenWithPublicAPI(cluster *mgmtv3.Cluster) (string, *bool, error) {
	// mock
	m.Mock.On("getRestConfig", cluster).Return(&rest.Config{}, nil)

	restConfig, err := m.getRestConfig(cluster)
	if err != nil {
		return "", nil, err
	}
	requiresTunnel := new(bool)
	restConfig.Dial = (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext

	// mock serviceToken
	serviceToken, err := "testtoken12345", nil

	if err != nil {
		*requiresTunnel = true
		var dnsError *net.DNSError
		if stderrors.As(err, &dnsError) && !dnsError.IsTemporary {
			return "", requiresTunnel, nil
		}
		var urlError *url.Error
		if stderrors.As(err, &urlError) && urlError.Timeout() {
			return "", requiresTunnel, nil
		}
		requiresTunnel = nil
	}
	return serviceToken, requiresTunnel, err
}

func (m *mockAksOperatorController) getRestConfig(cluster *mgmtv3.Cluster) (*rest.Config, error) {

	_mc_ret := m.Called(cluster)

	var _r0 *rest.Config

	if _rfn, ok := _mc_ret.Get(0).(func(*mgmtv3.Cluster) *rest.Config); ok {
		_r0 = _rfn(cluster)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(*rest.Config)
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

func getMockAksClusterConfig(filename string) (*unstructured.Unstructured, error) {
	var aksClusterConfig *unstructured.Unstructured

	// Read the embedded file
	bytes, err := testFs.ReadFile(filename); if err != nil {
		return aksClusterConfig, err
	}
	// Unmarshal json into an unstructured cluster config object
	err = json.Unmarshal(bytes, &aksClusterConfig); if err != nil {
		return aksClusterConfig, err
	}

	return aksClusterConfig, nil
}
