package digitalocean

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const digitalOceanCloudCredNameBase = "digitalOceanCloudCredential"

// CreateDigitalOceanCloudCredentials is a helper function that takes the ranger Client as a parameter and creates
// a Digital Ocean cloud credential, and returns the CloudCredential response
func CreateDigitalOceanCloudCredentials(rangerClient *ranger.Client) (*cloudcredentials.CloudCredential, error) {
	var digitalOceanCredentialConfig cloudcredentials.DigitalOceanCredentialConfig
	config.LoadConfig(cloudcredentials.DigitalOceanCredentialConfigurationFileKey, &digitalOceanCredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                         digitalOceanCloudCredNameBase,
		DigitalOceanCredentialConfig: &digitalOceanCredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
