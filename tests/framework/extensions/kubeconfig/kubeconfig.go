package kubeconfig

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeconfig generates a kubeconfig froma specific cluster, and returns it in the form of a *clientcmd.ClientConfig
func GetKubeconfig(client *ranger.Client, clusterID string) (*clientcmd.ClientConfig, error) {
	cluster, err := client.Management.Cluster.ByID(clusterID)
	if err != nil {
		return nil, err
	}

	kubeConfig, err := client.Management.Cluster.ActionGenerateKubeconfig(cluster)
	if err != nil {
		return nil, err
	}

	configBytes := []byte(kubeConfig.Config)

	clientConfig, err := clientcmd.NewClientConfigFromBytes(configBytes)
	if err != nil {
		return nil, err
	}

	return &clientConfig, nil
}
