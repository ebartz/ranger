package eks

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
)

// CreateEKSHostedCluster is a helper function that creates an EKS hosted cluster
func CreateEKSHostedCluster(client *ranger.Client, displayName, cloudCredentialID string, enableClusterAlerting, enableClusterMonitoring, enableNetworkPolicy, windowsPreferedCluster bool, labels map[string]string) (*management.Cluster, error) {
	eksHostCluster := eksHostClusterConfig(displayName, cloudCredentialID)
	cluster := &management.Cluster{
		DockerRootDir:           "/var/lib/docker",
		EKSConfig:               eksHostCluster,
		Name:                    displayName,
		EnableClusterAlerting:   enableClusterAlerting,
		EnableClusterMonitoring: enableClusterMonitoring,
		EnableNetworkPolicy:     &enableNetworkPolicy,
		Labels:                  labels,
		WindowsPreferedCluster:  windowsPreferedCluster,
	}

	clusterResp, err := client.Management.Cluster.Create(cluster)
	if err != nil {
		return nil, err
	}
	return clusterResp, err
}
