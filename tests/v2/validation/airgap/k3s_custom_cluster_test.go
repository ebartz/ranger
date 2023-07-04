package airgap

import (
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/corral"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/clusters/bundledclusters"
	"github.com/ranger/ranger/tests/framework/extensions/clusters/kubernetesversions"
	registryExtension "github.com/ranger/ranger/tests/framework/extensions/registries"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/ranger/ranger/tests/v2/validation/pipeline/rangerha/corralha"
	provisioning "github.com/ranger/ranger/tests/v2/validation/provisioning"
	"github.com/ranger/ranger/tests/v2/validation/registries"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AirGapK3SCustomClusterTestSuite struct {
	suite.Suite
	client             *ranger.Client
	session            *session.Session
	kubernetesVersions []string
	cnis               []string
	corralImage        string
	corralAutoCleanup  bool
	registryFQDN       string
	advancedOptions    provisioning.AdvancedOptions
}

func (a *AirGapK3SCustomClusterTestSuite) TearDownSuite() {
	a.session.Cleanup()
}

func (a *AirGapK3SCustomClusterTestSuite) SetupSuite() {
	testSession := session.NewSession()
	a.session = testSession

	clustersConfig := new(provisioning.Config)
	config.LoadConfig(provisioning.ConfigurationFileKey, clustersConfig)

	corralRangerHA := new(corralha.CorralRangerHA)
	config.LoadConfig(corralha.CorralRangerHAConfigConfigurationFileKey, corralRangerHA)

	registriesConfig := new(registries.Registries)
	config.LoadConfig(registries.RegistriesConfigKey, registriesConfig)

	a.kubernetesVersions = clustersConfig.K3SKubernetesVersions
	a.cnis = clustersConfig.CNIs
	a.advancedOptions = clustersConfig.AdvancedOptions

	client, err := ranger.NewClient("", testSession)
	require.NoError(a.T(), err)

	a.client = client
	listOfCorrals, err := corral.ListCorral()
	require.NoError(a.T(), err)

	corralConfig := corral.CorralConfigurations()

	err = corral.SetupCorralConfig(corralConfig.CorralConfigVars, corralConfig.CorralConfigUser, corralConfig.CorralSSHPath)
	require.NoError(a.T(), err)

	corralPackage := corral.CorralPackagesConfig()
	a.corralImage = corralPackage.CorralPackageImages[corralPackageAirgapCustomClusterName]
	a.corralAutoCleanup = corralPackage.HasCleanup

	_, corralExist := listOfCorrals[corralRangerHA.Name]
	if corralExist {
		bastionIP, err := corral.GetCorralEnvVar(corralRangerHA.Name, corralRegistryIP)
		require.NoError(a.T(), err)

		err = corral.UpdateCorralConfig(corralBastionIP, bastionIP)
		require.NoError(a.T(), err)

		registryFQDN, err := corral.GetCorralEnvVar(corralRangerHA.Name, corralRegistryFQDN)
		require.NoError(a.T(), err)
		logrus.Infof("registry fqdn is %s", registryFQDN)
		a.registryFQDN = registryFQDN
	} else {
		a.registryFQDN = registriesConfig.ExistingNoAuthRegistryURL
	}

}

func (a *AirGapK3SCustomClusterTestSuite) TestProvisioningK3SCustomCluster() {
	nodeRoles := map[int]string{
		1: "--etcd --controlplane --worker",
	}

	var name string
	for _, kubeVersion := range a.kubernetesVersions {
		name = "K3S Custom Cluster Kubernetes version: " + kubeVersion
		for _, cni := range a.cnis {
			name += " cni: " + cni
			a.Run(name, func() {
				clusterName := testProvisionAirgapRKE2K3SCustomCluster(a.T(), a.client, nodeRoles, a.corralImage, cni, kubeVersion, a.corralAutoCleanup, a.advancedOptions)
				passed, podErrors := registryExtension.CheckPodStatusImageSource(a.client, clusterName, a.registryFQDN)
				assert.Empty(a.T(), podErrors)
				assert.True(a.T(), passed)
			})
		}
	}
}

func (a *AirGapK3SCustomClusterTestSuite) TestProvisioningUpgradeK3SCustomCluster() {
	nodeRoles := map[int]string{
		1: "--etcd --controlplane --worker",
	}

	rke2Versions, err := kubernetesversions.ListK3SAllVersions(a.client)
	require.NoError(a.T(), err)

	numOfRKE2Versions := len(rke2Versions)
	// for this we will only have one custom cluster entry and one cni entry
	cni := a.cnis[0]
	kubeVersion := rke2Versions[numOfRKE2Versions-2]
	upgradeDefaultKubeVersion := rke2Versions[numOfRKE2Versions-1]

	clusterName := testProvisionAirgapRKE2K3SCustomCluster(a.T(), a.client, nodeRoles, a.corralImage, cni, kubeVersion, a.corralAutoCleanup, a.advancedOptions)
	clusterMeta, err := clusters.NewClusterMeta(a.client, clusterName)
	require.NoError(a.T(), err)
	require.NotNilf(a.T(), clusterMeta, "Couldn't get the cluster meta")

	initCluster, err := bundledclusters.NewWithClusterMeta(clusterMeta)
	require.NoError(a.T(), err)

	cluster, err := initCluster.Get(a.client)
	require.NoError(a.T(), err)

	updatedCluster, err := cluster.UpdateKubernetesVersion(a.client, &upgradeDefaultKubeVersion)
	require.NoError(a.T(), err)

	err = clusters.WaitClusterToBeUpgraded(a.client, clusterMeta.ID)
	require.NoError(a.T(), err)

	validateRKE2K3SKubernetesUpgrade(a.T(), updatedCluster, upgradeDefaultKubeVersion)

	passed, podErrors := registryExtension.CheckPodStatusImageSource(a.client, clusterName, a.registryFQDN)
	assert.Empty(a.T(), podErrors)
	assert.True(a.T(), passed)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAirGapCustomClusterK3SProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(AirGapK3SCustomClusterTestSuite))
}
