package harvester

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const harvesterCloudCredNameBase = "harvesterCloudCredential"

// CreateHarvesterCloudCredentials is a helper function that takes the ranger Client as a parameter and creates
// a harvester cloud credential, and returns the CloudCredential response
func CreateHarvesterCloudCredentials(rangerClient *ranger.Client) (*cloudcredentials.CloudCredential, error) {
	var harvesterCredentialConfig cloudcredentials.HarvesterCredentialConfig
	config.LoadConfig(cloudcredentials.HarvesterCredentialConfigurationFileKey, &harvesterCredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                      harvesterCloudCredNameBase,
		HarvesterCredentialConfig: &harvesterCredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
