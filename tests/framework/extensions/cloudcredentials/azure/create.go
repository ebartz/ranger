package azure

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const azureCloudCredNameBase = "azureOceanCloudCredential"

// CreateAzureCloudCredentials is a helper function that takes the ranger Client as a parameter and creates
// an Azure cloud credential, and returns the CloudCredential response
func CreateAzureCloudCredentials(rangerClient *ranger.Client) (*cloudcredentials.CloudCredential, error) {
	var azureCredentialConfig cloudcredentials.AzureCredentialConfig
	config.LoadConfig(cloudcredentials.AzureCredentialConfigurationFileKey, &azureCredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                  azureCloudCredNameBase,
		AzureCredentialConfig: &azureCredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
