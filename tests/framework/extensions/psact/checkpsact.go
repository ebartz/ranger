package psact

import (
	"fmt"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
)

// CheckPSACT checks to see if PSACT is enabled or not in the cluster.
func CheckPSACT(client *ranger.Client, clusterName string) error {
	clusterID, err := clusters.GetClusterIDByName(client, clusterName)
	if err != nil {
		return err
	}

	cluster, err := client.Management.Cluster.ByID(clusterID)
	if err != nil {
		return err
	}

	if cluster.DefaultPodSecurityAdmissionConfigurationTemplateName == "" {
		return fmt.Errorf("error: PSACT is not defined in this cluster!")
	}

	return nil
}
