package rke1

import (
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	nodepools "github.com/ranger/ranger/tests/framework/extensions/rke1/nodepools"
	"github.com/ranger/ranger/tests/framework/extensions/rke1/nodetemplates"
	"github.com/ranger/ranger/tests/framework/extensions/workloads/pods"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	namegen "github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	"github.com/ranger/ranger/tests/v2/validation/provisioning"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KdmChecksTestSuite struct {
	suite.Suite
	session                *session.Session
	client                 *ranger.Client
	ns                     string
	rke1kubernetesVersions []string
	cnis                   []string
	providers              []string
	nodesAndRoles          []nodepools.NodeRoles
	advancedOptions        provisioning.AdvancedOptions
}

const (
	defaultNamespace             = "default"
	ProvisioningSteveResouceType = "provisioning.cattle.io.cluster"
)

func (k *KdmChecksTestSuite) TearDownSuite() {
	k.session.Cleanup()
}

func (k *KdmChecksTestSuite) SetupSuite() {
	testSession := session.NewSession()
	k.session = testSession

	k.ns = defaultNamespace

	clustersConfig := new(provisioning.Config)
	config.LoadConfig(provisioning.ConfigurationFileKey, clustersConfig)

	k.rke1kubernetesVersions = clustersConfig.RKE1KubernetesVersions

	k.cnis = clustersConfig.CNIs
	k.providers = clustersConfig.Providers
	k.nodesAndRoles = clustersConfig.NodesAndRolesRKE1
	k.advancedOptions = clustersConfig.AdvancedOptions

	client, err := ranger.NewClient("", testSession)
	require.NoError(k.T(), err)

	k.client = client
}

func (k *KdmChecksTestSuite) TestRKE1K8sVersions() {
	logrus.Infof("checking for valid k8s versions..")
	require.GreaterOrEqual(k.T(), len(k.rke1kubernetesVersions), 1)
	// fetching all available k8s versions from ranger
	releasedK8sVersions, _ := kubernetesversions.ListRKE1AllVersions(k.client)
	logrus.Info("expected k8s versions : ", k.rke1kubernetesVersions)
	logrus.Info("k8s versions available on ranger server : ", releasedK8sVersions)
	for _, expectedK8sVersion := range k.rke1kubernetesVersions {
		require.Contains(k.T(), releasedK8sVersions, expectedK8sVersion)
	}
}

func (k *KdmChecksTestSuite) TestProvisioningSingleNodeRKE1Clusters() {
	require.GreaterOrEqual(k.T(), len(k.providers), 1)
	require.GreaterOrEqual(k.T(), len(k.cnis), 1)

	subSession := k.session.NewSession()
	defer subSession.Cleanup()

	client, err := k.client.WithSession(subSession)
	require.NoError(k.T(), err)

	providerName := k.providers[0]
	provider := CreateProvider(providerName)

	nodePools := []*management.NodePool{}
	nodePoolNames := []string{}
	clusterNames := []string{}
	clusterResps := []*management.Cluster{}

	for _, k8sVersion := range k.rke1kubernetesVersions {
		for _, cni := range k.cnis {
			logrus.Info("provisioning " + k8sVersion + " cluster..")
			nodeTemplate, err := provider.NodeTemplateFunc(client)
			require.NoError(k.T(), err)

			clusterResp, nodePool, nodePoolName, clusterName, err := k.provisionRKE1Cluster(client, provider, k.nodesAndRoles, k8sVersion, cni, nodeTemplate, k.advancedOptions)
			require.NoError(k.T(), err)

			nodePoolNames = append(nodePoolNames, nodePoolName)
			clusterNames = append(clusterNames, clusterName)
			clusterResps = append(clusterResps, clusterResp)
			nodePools = append(nodePools, nodePool)
		}
	}

	k.checkClustersReady(client, clusterResps, nodePools, clusterNames, nodePoolNames)

}

func (k *KdmChecksTestSuite) provisionRKE1Cluster(client *ranger.Client, provider Provider, nodesAndRoles []nodepools.NodeRoles, k8sVersion, cni string, nodeTemplate *nodetemplates.NodeTemplate, advancedOptions provisioning.AdvancedOptions) (*management.Cluster, *management.NodePool, string, string, error) {
	clusterName := namegen.AppendRandomString(provider.Name.String())

	cluster := clusters.NewRKE1ClusterConfig(clusterName, cni, k8sVersion, "", client, advancedOptions)
	clusterResp, err := clusters.CreateRKE1Cluster(client, cluster)
	require.NoError(k.T(), err)

	nodePool, err := nodepools.NodePoolSetup(client, nodesAndRoles, clusterResp.ID, nodeTemplate.ID)
	require.NoError(k.T(), err)

	nodePoolName := nodePool.Name

	return clusterResp, nodePool, nodePoolName, clusterName, nil
}

func (k *KdmChecksTestSuite) checkClustersReady(client *ranger.Client, clusterResps []*management.Cluster, nodePools []*management.NodePool, clusterNames []string, nodePoolNames []string) {
	for i, clusterResp := range clusterResps {
		opts := metav1.ListOptions{
			FieldSelector:  "metadata.name=" + clusterResp.ID,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		}

		logrus.Info("waiting for cluster ", clusterResp.Name, " with k8s version ", k.rke1kubernetesVersions[i], " to be up..")
		watchInterface, err := k.client.GetManagementWatchInterface(management.ClusterType, opts)
		require.NoError(k.T(), err)

		checkFunc := clusters.IsHostedProvisioningClusterReady

		err = wait.WatchWait(watchInterface, checkFunc)
		require.NoError(k.T(), err)
		assert.Equal(k.T(), clusterNames[i], clusterResp.Name)
		assert.Equal(k.T(), nodePoolNames[i], nodePools[i].Name)
		assert.Equal(k.T(), k.rke1kubernetesVersions[i], clusterResp.RangerKubernetesEngineConfig.Version)

		podResults, podErrors := pods.StatusPods(client, clusterResp.ID)
		assert.NotEmpty(k.T(), podResults)
		assert.Empty(k.T(), podErrors)

	}
}

func TestPostKdmOutOfBandReleaseChecks(t *testing.T) {
	suite.Run(t, new(KdmChecksTestSuite))
}
