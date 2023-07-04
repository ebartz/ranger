package linode

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const linodeCloudCredNameBase = "linodeCloudCredential"

// CreateLinodeCloudCredentials is a helper function that takes the ranger Client as a parameter and creates
// a Linode cloud credential, and returns the CloudCredential response
func CreateLinodeCloudCredentials(rangerClient *ranger.Client) (*cloudcredentials.CloudCredential, error) {
	var linodeCredentialConfig cloudcredentials.LinodeCredentialConfig
	config.LoadConfig(cloudcredentials.LinodeCredentialConfigurationFileKey, &linodeCredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                   linodeCloudCredNameBase,
		LinodeCredentialConfig: &linodeCredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
