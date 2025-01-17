package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/creasty/defaults"
	v3 "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/ranger/tests/framework/clients/k3d"
	rangerClient "github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/ranger/ranger/tests/framework/extensions/token"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	namegen "github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/retry"
)

var (
	agentTag         = os.Getenv("AGENT_TAG")
	masterAgentImage = "ranger/ranger-agent:" + agentTag
)

const (
	k3dClusterNameBasename = "k3d-cluster"
)

// setup for integration testing
func main() {
	rangerConfig := new(rangerClient.Config)

	user := &management.User{
		Username: "admin",
		Password: "admin",
	}

	logrus.Infof("Generating test config...")
	ipAddress := getOutboundIP()
	hostURL := fmt.Sprintf("%s:8443", ipAddress.String())
	token, err := token.GenerateUserToken(user, hostURL)
	if err != nil {
		logrus.Fatalf("error with generating admin token: %v", err)
	}

	clusterName := namegen.AppendRandomString(k3dClusterNameBasename)

	cleanup := true
	rangerConfig.AdminToken = token.Token
	rangerConfig.Host = hostURL
	rangerConfig.Cleanup = &cleanup
	rangerConfig.ClusterName = clusterName

	if err := defaults.Set(rangerConfig); err != nil {
		logrus.Fatalf("error with setting up config file: %v", err)
	}

	config.WriteConfig(rangerClient.ConfigurationFileKey, rangerConfig)

	logrus.Infof("Setting up K3D downstream cluster...")
	testSession := session.NewSession()

	client, err := rangerClient.NewClient("", testSession)
	if err != nil {
		logrus.Fatalf("error creating admin client: %v", err)
	}

	agentSetting := &v3.Setting{}

	agentSettingResp, err := client.Steve.SteveType("management.cattle.io.setting").ByID("agent-image")
	if err != nil {
		logrus.Fatalf("error get agent-image setting: %v", err)
	}

	err = v1.ConvertToK8sType(agentSettingResp.JSONResp, agentSetting)
	if err != nil {
		logrus.Fatalf("error converting to k8s type: %v", err)
	}

	agentSetting.Value = masterAgentImage

	_, err = client.Steve.SteveType("management.cattle.io.setting").Update(agentSettingResp, agentSetting)
	if err != nil {
		logrus.Fatalf("error updating agent-image setting: %v", err)
	}
	logrus.Infof("Updated agent-image setting to %s", agentSetting.Value)

	// docker is sometimes unable to take the xtables lock to set up networking.
	// See this issue https://github.com/weaveworks/scope/issues/2308 which describes similar symptoms,
	// and points to https://github.com/moby/moby/issues/10218 (still open as of April 21, 2023) as a possible root cause.
	err = retry.OnError(retry.DefaultBackoff, func(err error) bool {
		return strings.Contains(err.Error(), "iptables: Resource temporarily unavailable")
	}, func() error {
		_, err = k3d.CreateAndImportK3DCluster(client, clusterName, masterAgentImage, "", 1, 0, true)
		return err
	})
	if err != nil {
		logrus.Fatalf("error creating and importing a k3d cluster: %v", err)
	}
}

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
