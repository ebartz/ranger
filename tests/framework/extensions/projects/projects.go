package projects

import (
	"github.com/ranger/norman/types"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
)

// GetProjectByName is a helper function that returns the project by name in a specific cluster.
func GetProjectByName(client *ranger.Client, clusterID, projectName string) (*management.Project, error) {
	var project *management.Project

	adminClient, err := ranger.NewClient(client.RangerConfig.AdminToken, client.Session)
	if err != nil {
		return project, err
	}

	projectsList, err := adminClient.Management.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": clusterID,
		},
	})
	if err != nil {
		return project, err
	}

	for i, p := range projectsList.Data {
		if p.Name == projectName {
			project = &projectsList.Data[i]
			break
		}
	}

	return project, nil
}

// GetProjectList is a helper function that returns all the project in a specific cluster
func GetProjectList(client *ranger.Client, clusterID string) (*management.ProjectCollection, error) {
	var projectsList *management.ProjectCollection

	projectsList, err := client.Management.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": clusterID,
		},
	})
	if err != nil {
		return projectsList, err
	}

	return projectsList, nil
}
