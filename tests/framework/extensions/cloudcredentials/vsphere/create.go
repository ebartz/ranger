package vsphere

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/cloudcredentials"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

const vmwarevsphereCloudCredNameBase = "vmwarevsphereCloudCredential"

// CreateVsphereCloudCredentials is a helper function that takes the ranger Client as a parameter and creates
// an AWS cloud credential, and returns the CloudCredential response
func CreateVsphereCloudCredentials(rangerClient *ranger.Client) (*cloudcredentials.CloudCredential, error) {
	var vmwarevsphereCredentialConfig cloudcredentials.VmwarevsphereCredentialConfig
	config.LoadConfig(cloudcredentials.VmwarevsphereCredentialConfigurationFileKey, &vmwarevsphereCredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                vmwarevsphereCloudCredNameBase,
		VmwareVsphereConfig: &vmwarevsphereCredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rangerClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
