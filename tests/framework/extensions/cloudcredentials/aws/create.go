package aws

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const awsCloudCredNameBase = "awsCloudCredential"

// CreateAWSCloudCredentials is a helper function that takes the ranger Client as a parameter and creates
// an AWS cloud credential, and returns the CloudCredential response
func CreateAWSCloudCredentials(rangerClient *ranger.Client) (*cloudcredentials.CloudCredential, error) {
	var amazonEC2CredentialConfig cloudcredentials.AmazonEC2CredentialConfig
	config.LoadConfig(cloudcredentials.AmazonEC2CredentialConfigurationFileKey, &amazonEC2CredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                      awsCloudCredNameBase,
		AmazonEC2CredentialConfig: &amazonEC2CredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
