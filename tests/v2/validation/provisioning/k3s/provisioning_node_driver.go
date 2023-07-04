package k3s

import (
	"context"
	"fmt"
	"testing"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	"github.com/ranger/ranger/tests/framework/extensions/machinepools"
	nodestat "github.com/ranger/ranger/tests/framework/extensions/nodes"
	"github.com/ranger/ranger/tests/framework/extensions/pipeline"
	psadeploy "github.com/ranger/ranger/tests/framework/extensions/psact"
	"github.com/ranger/ranger/tests/framework/extensions/workloads/pods"
	"github.com/ranger/ranger/tests/framework/pkg/environmentflag"
	namegen "github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	provisioning "github.com/ranger/ranger/tests/v2/validation/provisioning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "fleet-default"
)

func TestProvisioningK3SCluster(t *testing.T, client *ranger.Client, provider Provider, nodesAndRoles []machinepools.NodeRoles, kubeVersion string, psact string, advancedOptions provisioning.AdvancedOptions) *v1.SteveAPIObject {
	cloudCredential, err := provider.CloudCredFunc(client)
	require.NoError(t, err)

	clusterName := namegen.AppendRandomString(provider.Name.String())
	generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
	machinePoolConfig := provider.MachinePoolFunc(generatedPoolName, namespace)

	machineConfigResp, err := client.Steve.SteveType(provider.MachineConfigPoolResourceSteveType).Create(machinePoolConfig)
	require.NoError(t, err)

	machinePools := machinepools.RKEMachinePoolSetup(nodesAndRoles, machineConfigResp)

	cluster := clusters.NewK3SRKE2ClusterConfig(clusterName, namespace, "", cloudCredential.ID, kubeVersion, psact, machinePools, advancedOptions)

	clusterResp, err := clusters.CreateK3SRKE2Cluster(client, cluster)
	require.NoError(t, err)

	if client.Flags.GetValue(environmentflag.UpdateClusterName) {
		pipeline.UpdateConfigClusterName(clusterName)
	}

	adminClient, err := ranger.NewClient(client.RangerConfig.AdminToken, client.Session)
	require.NoError(t, err)
	kubeProvisioningClient, err := adminClient.GetKubeAPIProvisioningClient()
	require.NoError(t, err)

	result, err := kubeProvisioningClient.Clusters(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + clusterName,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	require.NoError(t, err)

	checkFunc := clusters.IsProvisioningClusterReady

	err = wait.WatchWait(result, checkFunc)
	assert.NoError(t, err)
	assert.Equal(t, clusterName, clusterResp.ObjectMeta.Name)
	assert.Equal(t, kubeVersion, cluster.Spec.KubernetesVersion)

	clusterIDName, err := clusters.GetClusterIDByName(adminClient, clusterName)
	assert.NoError(t, err)

	err = nodestat.IsNodeReady(client, clusterIDName)
	require.NoError(t, err)

	clusterToken, err := clusters.CheckServiceAccountTokenSecret(client, clusterName)
	require.NoError(t, err)
	assert.NotEmpty(t, clusterToken)

	if psact == string(provisioning.RangerPrivileged) || psact == string(provisioning.RangerRestricted) {
		err = psadeploy.CheckPSACT(client, clusterName)
		require.NoError(t, err)

		_, err = psadeploy.CreateNginxDeployment(client, clusterIDName, psact)
		require.NoError(t, err)
	}

	podResults, podErrors := pods.StatusPods(client, clusterIDName)
	assert.NotEmpty(t, podResults)
	assert.Empty(t, podErrors)

	return clusterResp
}
