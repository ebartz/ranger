package nodetemplates

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/rke1/nodetemplates"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const harvesterNodeTemplateNameBase = "harvesterNodeConfig"

// CreateHarvesterNodeTemplate is a helper function that takes the ranger Client as a parameter and creates
// an Harvester node template and returns the NodeTemplate response
func CreateHarvesterNodeTemplate(rangerClient *ranger.Client) (*nodetemplates.NodeTemplate, error) {
	var harvesterNodeTemplateConfig nodetemplates.HarvesterNodeTemplateConfig
	config.LoadConfig(nodetemplates.HarvesterNodeTemplateConfigurationFileKey, &harvesterNodeTemplateConfig)

	nodeTemplate := nodetemplates.NodeTemplate{
		EngineInstallURL:            "https://releases.ranger.com/install-docker/24.0.sh",
		Name:                        harvesterNodeTemplateNameBase,
		HarvesterNodeTemplateConfig: &harvesterNodeTemplateConfig,
	}

	resp := &nodetemplates.NodeTemplate{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.NodeTemplateType, nodeTemplate, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
