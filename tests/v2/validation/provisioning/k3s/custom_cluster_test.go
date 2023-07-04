package k3s

import (
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/ranger/ranger/tests/framework/extensions/machinepools"
	"github.com/ranger/ranger/tests/framework/extensions/users"
	password "github.com/ranger/ranger/tests/framework/extensions/users/passwordgenerator"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	namegen "github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	provisioning "github.com/ranger/ranger/tests/v2/validation/provisioning"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CustomClusterProvisioningTestSuite struct {
	suite.Suite
	client             *ranger.Client
	session            *session.Session
	standardUserClient *ranger.Client
	kubernetesVersions []string
	nodeProviders      []string
	psact              string
	hardened           bool
	advancedOptions    provisioning.AdvancedOptions
}

func (c *CustomClusterProvisioningTestSuite) TearDownSuite() {
	c.session.Cleanup()
}

func (c *CustomClusterProvisioningTestSuite) SetupSuite() {
	testSession := session.NewSession()
	c.session = testSession

	clustersConfig := new(provisioning.Config)
	config.LoadConfig(provisioning.ConfigurationFileKey, clustersConfig)

	c.kubernetesVersions = clustersConfig.K3SKubernetesVersions
	c.nodeProviders = clustersConfig.NodeProviders
	c.psact = clustersConfig.PSACT
	c.hardened = clustersConfig.Hardened
	c.advancedOptions = clustersConfig.AdvancedOptions

	client, err := ranger.NewClient("", testSession)
	require.NoError(c.T(), err)

	c.client = client

	c.kubernetesVersions, err = kubernetesversions.Default(c.client, clusters.K3SClusterType.String(), c.kubernetesVersions)
	require.NoError(c.T(), err)

	enabled := true
	var testuser = namegen.AppendRandomString("testuser-")
	var testpassword = password.GenerateUserPassword("testpass-")
	user := &management.User{
		Username: testuser,
		Password: testpassword,
		Name:     testuser,
		Enabled:  &enabled,
	}

	newUser, err := users.CreateUserWithRole(client, user, "user")
	require.NoError(c.T(), err)

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	require.NoError(c.T(), err)

	c.standardUserClient = standardUserClient
}

func (c *CustomClusterProvisioningTestSuite) TestProvisioningK3SCustomCluster() {
	nodeRolesAll := []machinepools.NodeRoles{provisioning.AllRolesPool}
	nodeRolesShared := []machinepools.NodeRoles{provisioning.EtcdControlPlanePool, provisioning.WorkerPool}
	nodeRolesDedicated := []machinepools.NodeRoles{provisioning.EtcdPool, provisioning.ControlPlanePool, provisioning.WorkerPool}

	tests := []struct {
		name      string
		client    *ranger.Client
		nodeRoles []machinepools.NodeRoles
		psact     string
	}{
		{"1 Node all roles " + provisioning.AdminClientName.String(), c.client, nodeRolesAll, c.psact},
		{"1 Node all roles " + provisioning.StandardClientName.String(), c.standardUserClient, nodeRolesAll, c.psact},
		{"2 nodes - etcd/cp roles per 1 node " + provisioning.AdminClientName.String(), c.client, nodeRolesShared, c.psact},
		{"2 nodes - etcd/cp roles per 1 node " + provisioning.StandardClientName.String(), c.standardUserClient, nodeRolesShared, c.psact},
		{"3 nodes - 1 role per node " + provisioning.AdminClientName.String(), c.client, nodeRolesDedicated, c.psact},
		{"3 nodes - 1 role per node " + provisioning.StandardClientName.String(), c.standardUserClient, nodeRolesDedicated, c.psact},
	}
	var name string
	for _, tt := range tests {
		testSession := session.NewSession()
		defer testSession.Cleanup()

		client, err := tt.client.WithSession(testSession)
		require.NoError(c.T(), err)

		for _, nodeProviderName := range c.nodeProviders {
			externalNodeProvider := provisioning.ExternalNodeProviderSetup(nodeProviderName)
			providerName := " Node Provider: " + nodeProviderName
			for _, kubeVersion := range c.kubernetesVersions {
				name = tt.name + providerName + " Kubernetes version: " + kubeVersion
				c.Run(name, func() {
					TestProvisioningK3SCustomCluster(c.T(), client, externalNodeProvider, tt.nodeRoles, kubeVersion, c.hardened, tt.psact, c.advancedOptions)
				})
			}
		}
	}
}

func (c *CustomClusterProvisioningTestSuite) TestProvisioningK3SCustomClusterDynamicInput() {
	clustersConfig := new(provisioning.Config)
	config.LoadConfig(provisioning.ConfigurationFileKey, clustersConfig)
	nodesAndRoles := clustersConfig.NodesAndRoles

	if len(nodesAndRoles) == 0 {
		c.T().Skip()
	}

	tests := []struct {
		name   string
		client *ranger.Client
		psact  string
	}{
		{provisioning.AdminClientName.String(), c.client, c.psact},
		{provisioning.StandardClientName.String(), c.standardUserClient, c.psact},
	}
	var name string
	for _, tt := range tests {
		testSession := session.NewSession()
		defer testSession.Cleanup()

		client, err := tt.client.WithSession(testSession)
		require.NoError(c.T(), err)

		for _, nodeProviderName := range c.nodeProviders {
			externalNodeProvider := provisioning.ExternalNodeProviderSetup(nodeProviderName)
			providerName := " Node Provider: " + nodeProviderName
			for _, kubeVersion := range c.kubernetesVersions {
				name = tt.name + providerName + " Kubernetes version: " + kubeVersion
				c.Run(name, func() {
					TestProvisioningK3SCustomCluster(c.T(), client, externalNodeProvider, nodesAndRoles, kubeVersion, c.hardened, tt.psact, c.advancedOptions)
				})
			}
		}
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCustomClusterK3SProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(CustomClusterProvisioningTestSuite))
}
