package pipeline

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/pkg/config"
)

// UpdateConfig is function that updates the cattle config's cluster name field which is
// the child of the ranger key in the cattle configuration.
func UpdateConfigClusterName(clusterName string) {
	rangerConfig := new(ranger.Config)
	config.LoadAndUpdateConfig(ranger.ConfigurationFileKey, rangerConfig, func() {
		rangerConfig.ClusterName = clusterName
	})
}
