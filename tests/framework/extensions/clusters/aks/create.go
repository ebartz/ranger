package aks

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
)

// CreateAKSHostedCluster is a helper function that creates an AKS hosted cluster.
func CreateAKSHostedCluster(client *ranger.Client, displayName, cloudCredentialID string, enableClusterAlerting, enableClusterMonitoring, enableNetworkPolicy, windowsPreferedCluster bool, labels map[string]string) (*management.Cluster, error) {
	aksHostCluster := AKSHostClusterConfig(displayName, cloudCredentialID)
	cluster := &management.Cluster{
		AKSConfig:               aksHostCluster,
		DockerRootDir:           "/var/lib/docker",
		EnableClusterAlerting:   enableClusterAlerting,
		EnableClusterMonitoring: enableClusterMonitoring,
		EnableNetworkPolicy:     &enableNetworkPolicy,
		Labels:                  labels,
		Name:                    displayName,
		WindowsPreferedCluster:  windowsPreferedCluster,
	}

	clusterResp, err := client.Management.Cluster.Create(cluster)
	if err != nil {
		return nil, err
	}

	return clusterResp, err
}
