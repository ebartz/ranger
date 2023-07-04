package nodetemplates

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/rke1/nodetemplates"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const vmwarevsphereNodeTemplateNameBase = "vmwarevsphereNodeConfig"

// CreateVSphereNodeTemplate is a helper function that takes the ranger Client as a parameter and creates
// an VSphere node template and returns the NodeTemplate response
func CreateVSphereNodeTemplate(rangerClient *ranger.Client) (*nodetemplates.NodeTemplate, error) {
	var vmwarevsphereNodeTemplateConfig nodetemplates.VmwareVsphereNodeTemplateConfig
	config.LoadConfig(nodetemplates.VmwareVsphereNodeTemplateConfigurationFileKey, &vmwarevsphereNodeTemplateConfig)

	nodeTemplate := nodetemplates.NodeTemplate{
		EngineInstallURL:                "https://releases.ranger.com/install-docker/20.10.sh",
		Name:                            vmwarevsphereNodeTemplateNameBase,
		VmwareVsphereNodeTemplateConfig: &vmwarevsphereNodeTemplateConfig,
	}

	resp := &nodetemplates.NodeTemplate{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.NodeTemplateType, nodeTemplate, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
