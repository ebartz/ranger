package nodetemplates

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/rke1/nodetemplates"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const awsEC2NodeTemplateNameBase = "awsNodeConfig"

// CreateAWSNodeTemplate is a helper function that takes the ranger Client as a parameter and creates
// an AWS node template and returns the NodeTemplate response
func CreateAWSNodeTemplate(rangerClient *ranger.Client) (*nodetemplates.NodeTemplate, error) {
	var amazonEC2NodeTemplateConfig nodetemplates.AmazonEC2NodeTemplateConfig
	config.LoadConfig(nodetemplates.AmazonEC2NodeTemplateConfigurationFileKey, &amazonEC2NodeTemplateConfig)

	nodeTemplate := nodetemplates.NodeTemplate{
		EngineInstallURL:            "https://releases.ranger.com/install-docker/24.0.sh",
		Name:                        awsEC2NodeTemplateNameBase,
		AmazonEC2NodeTemplateConfig: &amazonEC2NodeTemplateConfig,
	}

	resp := &nodetemplates.NodeTemplate{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.NodeTemplateType, nodeTemplate, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
