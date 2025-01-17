package provisioning

import (
	"fmt"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/nodes/ec2"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	"github.com/ranger/ranger/tests/framework/pkg/nodes"
)

const (
	ec2NodeProviderName = "ec2"
	fromConfig          = "config"
)

type NodeCreationFunc func(client *ranger.Client, rolesPerPool []string, quantityPerPool []int32) (nodes []*nodes.Node, err error)

type ExternalNodeProvider struct {
	Name             string
	NodeCreationFunc NodeCreationFunc
}

// ExternalNodeProviderSetup is a helper function that setups an ExternalNodeProvider object is a wrapper
// for the specific outside node provider node creator function
func ExternalNodeProviderSetup(providerType string) ExternalNodeProvider {
	switch providerType {
	case ec2NodeProviderName:
		return ExternalNodeProvider{
			Name:             providerType,
			NodeCreationFunc: ec2.CreateNodes,
		}
	case fromConfig:
		return ExternalNodeProvider{
			Name: providerType,
			NodeCreationFunc: func(client *ranger.Client, rolesPerPool []string, quantityPerPool []int32) (nodesList []*nodes.Node, err error) {
				var nodeConfig nodes.ExternalNodeConfig
				config.LoadConfig(nodes.ExternalNodeConfigConfigurationFileKey, &nodeConfig)

				nodesList = nodeConfig.Nodes[-1]

				for _, node := range nodesList {
					sshKey, err := nodes.GetSSHKey(node.SSHKeyName)
					if err != nil {
						return nil, err
					}

					node.SSHKey = sshKey
				}
				return nodesList, nil
			},
		}
	default:
		panic(fmt.Sprintf("Node Provider:%v not found", providerType))
	}

}
