package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ranger/norman/types"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	"github.com/ranger/ranger/tests/framework/extensions/token"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
)

type wrappedConfig struct {
	Configuration *ranger.Config `yaml:"ranger"`
}

var (
	adminPassword = os.Getenv("ADMIN_PASSWORD")
	host          = os.Getenv("HA_HOST")

	configFileName = config.ConfigFileName("cattle-config.yaml")
)

func main() {
	rangerConfig := new(ranger.Config)
	rangerConfig.Host = host
	isCleanupEnabled := true
	rangerConfig.Cleanup = &isCleanupEnabled

	adminUser := &management.User{
		Username: "admin",
		Password: adminPassword,
	}

	//create admin token
	var adminToken *management.Token
	err := kwait.Poll(500*time.Millisecond, 5*time.Minute, func() (done bool, err error) {
		adminToken, err = token.GenerateUserToken(adminUser, host)
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		logrus.Errorf("error creating admin token", err)
	}
	rangerConfig.AdminToken = adminToken.Token

	//create config file
	configWrapped := &wrappedConfig{
		Configuration: rangerConfig,
	}
	configData, err := yaml.Marshal(configWrapped)
	if err != nil {
		logrus.Errorf("error marshaling", err)
	}
	err = configFileName.NewFile(configData)
	if err != nil {
		logrus.Fatalf("error writing yaml", err)
	}
	err = configFileName.SetEnvironmentKey()
	if err != nil {
		logrus.Fatalf("error while setting environment path", err)
	}

	session := session.NewSession()
	client, err := ranger.NewClient("", session)
	if err != nil {
		logrus.Errorf("error creating client", err)
	}

	clusterList, err := client.Management.Cluster.List(&types.ListOpts{})
	if err != nil {
		logrus.Errorf("error getting cluster list", err)
	}

	for _, c := range clusterList.Data {
		isLocalCluster := c.ID == "local"
		if !isLocalCluster {
			opts := metav1.ListOptions{
				FieldSelector:  "metadata.name=" + c.ID,
				TimeoutSeconds: &defaults.WatchTimeoutSeconds,
			}

			err := client.Management.Cluster.Delete(&c)
			if err != nil {
				logrus.Errorf("error delete cluster call: %v", err)
			}

			watchInterface, err := client.GetManagementWatchInterface(management.ClusterType, opts)
			if err != nil {
				logrus.Errorf("error while getting the watch interface: %v", err)
			}

			wait.WatchWait(watchInterface, func(event watch.Event) (ready bool, err error) {
				if event.Type == watch.Error {
					return false, fmt.Errorf("there was an error deleting cluster")
				} else if event.Type == watch.Deleted {
					return true, nil
				}
				return false, nil
			})
		}
	}
}
