package provisioning

import (
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials/azure"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/clusters/aks"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	nodestat "github.com/ranger/ranger/tests/framework/extensions/nodes"
	"github.com/ranger/ranger/tests/framework/extensions/pipeline"
	"github.com/ranger/ranger/tests/framework/extensions/users"
	password "github.com/ranger/ranger/tests/framework/extensions/users/passwordgenerator"
	"github.com/ranger/ranger/tests/framework/extensions/workloads/pods"
	"github.com/ranger/ranger/tests/framework/pkg/environmentflag"
	namegen "github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	"github.com/ranger/ranger/tests/v2/validation/provisioning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HostedAKSClusterProvisioningTestSuite struct {
	suite.Suite
	client             *ranger.Client
	session            *session.Session
	standardUserClient *ranger.Client
	cluster            *management.Cluster
}

func (h *HostedAKSClusterProvisioningTestSuite) TearDownSuite() {
	h.session.Cleanup()
}

func (h *HostedAKSClusterProvisioningTestSuite) SetupSuite() {
	testSession := session.NewSession()
	h.session = testSession

	client, err := ranger.NewClient("", testSession)
	require.NoError(h.T(), err)

	h.client = client

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
	require.NoError(h.T(), err)

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	require.NoError(h.T(), err)

	h.standardUserClient = standardUserClient
}

func (h *HostedAKSClusterProvisioningTestSuite) TestProvisioningHostedAKS() {
	tests := []struct {
		name   string
		client *ranger.Client
	}{
		{provisioning.AdminClientName.String(), h.client},
		{provisioning.StandardClientName.String(), h.standardUserClient},
	}

	for _, tt := range tests {
		subSession := h.session.NewSession()
		defer subSession.Cleanup()

		client, err := tt.client.WithSession(subSession)
		require.NoError(h.T(), err)

		cloudCredential, err := azure.CreateAzureCloudCredentials(client)
		require.NoError(h.T(), err)

		scaleName := "scaling " + tt.name
		h.Run(tt.name, func() {
			cluster, err := h.testProvisioningHostedAKSCluster(client, cloudCredential)
			require.NoError(h.T(), err)

			h.cluster = cluster
		})

		h.Run(scaleName, func() {
			h.testScalingAKSNodePools(client, h.cluster, cloudCredential)
		})

		h.cluster = nil
	}
}

func (h *HostedAKSClusterProvisioningTestSuite) testProvisioningHostedAKSCluster(rangerClient *ranger.Client, cloudCredential *cloudcredentials.CloudCredential) (*management.Cluster, error) {
	clusterName := namegen.AppendRandomString("akshostcluster")
	clusterResp, err := aks.CreateAKSHostedCluster(rangerClient, clusterName, cloudCredential.ID, false, false, false, false, map[string]string{})
	require.NoError(h.T(), err)

	if h.client.Flags.GetValue(environmentflag.UpdateClusterName) {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	opts := metav1.ListOptions{
		FieldSelector:  "metadata.name=" + clusterResp.ID,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	}
	watchInterface, err := h.client.GetManagementWatchInterface(management.ClusterType, opts)
	require.NoError(h.T(), err)

	checkFunc := clusters.IsHostedProvisioningClusterReady

	err = wait.WatchWait(watchInterface, checkFunc)
	require.NoError(h.T(), err)
	assert.Equal(h.T(), clusterName, clusterResp.Name)

	clusterToken, err := clusters.CheckServiceAccountTokenSecret(rangerClient, clusterName)
	require.NoError(h.T(), err)
	assert.NotEmpty(h.T(), clusterToken)

	err = nodestat.IsNodeReady(rangerClient, clusterResp.ID)
	require.NoError(h.T(), err)

	podResults, podErrors := pods.StatusPods(rangerClient, clusterResp.ID)
	assert.NotEmpty(h.T(), podResults)
	assert.Empty(h.T(), podErrors)

	return clusterResp, nil
}

func (h *HostedAKSClusterProvisioningTestSuite) testScalingAKSNodePools(rangerClient *ranger.Client, cluster *management.Cluster, cloudCredential *cloudcredentials.CloudCredential) {
	if cluster == nil {
		cluster, err := h.testProvisioningHostedAKSCluster(rangerClient, cloudCredential)
		require.NoError(h.T(), err)

		updatedCluster, err := ScalingAKSNodePools(rangerClient, cluster, cluster.Name, cloudCredential)
		require.NoError(h.T(), err)
		assert.Equal(h.T(), cluster.Name, updatedCluster.Name)

	} else {
		updatedCluster, err := ScalingAKSNodePools(rangerClient, cluster, cluster.Name, cloudCredential)
		require.NoError(h.T(), err)
		assert.Equal(h.T(), cluster.Name, updatedCluster.Name)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHostedAKSClusterProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(HostedAKSClusterProvisioningTestSuite))
}
