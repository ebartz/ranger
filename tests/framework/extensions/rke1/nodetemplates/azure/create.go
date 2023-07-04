package nodetemplates

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/rke1/nodetemplates"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const azureNodeTemplateNameBase = "azureNodeConfig"

// CreateAzureNodeTemplate is a helper function that takes the ranger Client as a parameter and creates
// an Azure node template and returns the NodeTemplate response
func CreateAzureNodeTemplate(rangerClient *ranger.Client) (*nodetemplates.NodeTemplate, error) {
	var azureNodeTemplateConfig nodetemplates.AzureNodeTemplateConfig
	config.LoadConfig(nodetemplates.AzureNodeTemplateConfigurationFileKey, &azureNodeTemplateConfig)

	nodeTemplate := nodetemplates.NodeTemplate{
		EngineInstallURL:        "https://releases.ranger.com/install-docker/23.0.sh",
		Name:                    azureNodeTemplateNameBase,
		AzureNodeTemplateConfig: &azureNodeTemplateConfig,
	}

	resp := &nodetemplates.NodeTemplate{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.NodeTemplateType, nodeTemplate, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
