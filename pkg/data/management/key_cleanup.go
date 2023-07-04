package management

import (
	"github.com/ranger/ranger/pkg/types/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func sshKeyCleanup(management *config.ManagementContext) error {
	nodeClient := management.Management.Nodes("")
	nodes, err := nodeClient.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range nodes.Items {
		if node.Status.NodeConfig != nil && node.Status.NodeConfig.SSHKey != "" {
			node.Status.NodeConfig.SSHKey = ""
			_, err = nodeClient.Update(&node)
			if err != nil {
				return err
			}
		}
	}

	clusterClient := management.Management.Clusters("")
	clusters, err := clusterClient.List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cluster := range clusters.Items {
		if cluster.Status.AppliedSpec.RangerKubernetesEngineConfig == nil {
			continue
		}
		pruned := false
		for i, node := range cluster.Status.AppliedSpec.RangerKubernetesEngineConfig.Nodes {
			if node.SSHKey != "" {
				cluster.Status.AppliedSpec.RangerKubernetesEngineConfig.Nodes[i].SSHKey = ""
				pruned = true
			}
		}
		if pruned {
			_, err = clusterClient.Update(&cluster)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
