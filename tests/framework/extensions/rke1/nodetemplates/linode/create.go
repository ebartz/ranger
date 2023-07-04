package nodetemplates

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/rke1/nodetemplates"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const linodeNodeTemplateNameBase = "linodeNodeConfig"

// CreateLinodeNodeTemplate is a helper function that takes the ranger Client as a parameter and creates
// an Linode node template and returns the NodeTemplate response
func CreateLinodeNodeTemplate(rangerClient *ranger.Client) (*nodetemplates.NodeTemplate, error) {
	var linodeNodeTemplateConfig nodetemplates.LinodeNodeTemplateConfig
	config.LoadConfig(nodetemplates.LinodeNodeTemplateConfigurationFileKey, &linodeNodeTemplateConfig)

	nodeTemplate := nodetemplates.NodeTemplate{
		EngineInstallURL:         "https://releases.ranger.com/install-docker/24.0.sh",
		Name:                     linodeNodeTemplateNameBase,
		LinodeNodeTemplateConfig: &linodeNodeTemplateConfig,
	}

	resp := &nodetemplates.NodeTemplate{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.NodeTemplateType, nodeTemplate, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
