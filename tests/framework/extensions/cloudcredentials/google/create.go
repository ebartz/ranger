package google

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const googleCloudCredNameBase = "googleCloudCredNameBase"

// CreateGoogleCloudCredentials is a helper function that takes the ranger Client as a parameter and creates
// a Google cloud credential, and returns the CloudCredential response
func CreateGoogleCloudCredentials(rangerClient *ranger.Client) (*cloudcredentials.CloudCredential, error) {
	var googleCredentialConfig cloudcredentials.GoogleCredentialConfig
	config.LoadConfig(cloudcredentials.GoogleCredentialConfigurationFileKey, &googleCredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                   googleCloudCredNameBase,
		GoogleCredentialConfig: &googleCredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
