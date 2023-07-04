package main

import (
	"context"
	"os"
	"time"

	"github.com/ranger/ranger/pkg/api/scheme"
	"github.com/ranger/ranger/tests/framework/clients/corral"
	"github.com/ranger/ranger/tests/framework/clients/dynamic"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	"github.com/ranger/ranger/tests/framework/extensions/kubeapi/workloads/deployments"
	"github.com/ranger/ranger/tests/framework/extensions/pipeline"
	nodepools "github.com/ranger/ranger/tests/framework/extensions/rke1/nodepools"
	aws "github.com/ranger/ranger/tests/framework/extensions/rke1/nodetemplates/aws"
	"github.com/ranger/ranger/tests/framework/extensions/token"
	"github.com/ranger/ranger/tests/framework/extensions/unstructured"
	"github.com/ranger/ranger/tests/framework/extensions/users"
	passwordgenerator "github.com/ranger/ranger/tests/framework/extensions/users/passwordgenerator"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	namegen "github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	"github.com/ranger/ranger/tests/v2/validation/provisioning"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	appv1 "k8s.io/api/apps/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

const (
	corralName       = "rangertestcoverage"
	rangerTestImage = "rangertest/ranger:v2.7-head"
	namespace        = "cattle-system"
	deploymentName   = "ranger"
	clusterName      = "local"
	// The json/yaml config key for the corral package to be build ..
	userClusterConfigsConfigurationFileKey = "userClusterConfig"
)

type UserClusterConfig struct {
	Token         string   `json:"token" yaml:"token"`
	Username      string   `json:"username" yaml:"username"`
	Clusters      []string `json:"clusters" yaml:"clusters"`
	AdminPassword string   `json:"adminPassword" yaml:"adminPassword"`
}

// setup for code coverage testing and reporting
func main() {
	rangerConfig := new(ranger.Config)
	config.LoadConfig(ranger.ConfigurationFileKey, rangerConfig)

	kubeconfig, err := getRangerKubeconfig()
	if err != nil {
		logrus.Fatalf("error with getting kube config using corral: %v", err)
	}

	password, err := corral.GetCorralEnvVar(corralName, "bootstrap_password")
	if err != nil {
		logrus.Fatalf("error getting password %v", err)
	}

	// update deployment
	err = updateRangerDeployment(kubeconfig)
	if err != nil {
		logrus.Fatalf("error updating ranger deployment: %v", err)
	}

	token, err := createAdminToken(password, rangerConfig)
	if err != nil {
		logrus.Fatalf("error with generating admin token: %v", err)
	}

	rangerConfig.AdminToken = token
	config.UpdateConfig(ranger.ConfigurationFileKey, rangerConfig)
	//provision clusters for test
	session := session.NewSession()
	clustersConfig := new(provisioning.Config)
	config.LoadConfig(provisioning.ConfigurationFileKey, clustersConfig)
	kubernetesVersions := clustersConfig.RKE1KubernetesVersions
	cnis := clustersConfig.CNIs
	nodesAndRoles := clustersConfig.NodesAndRolesRKE1
	advancedOptions := clustersConfig.AdvancedOptions

	client, err := ranger.NewClient("", session)
	if err != nil {
		logrus.Fatalf("error creating admin client: %v", err)
	}

	err = pipeline.PostRangerInstall(client, adminPassword)
	if err != nil {
		logrus.Errorf("error during post ranger install: %v", err)
	}

	enabled := true
	var testuser = namegen.AppendRandomString("testuser-")
	var testpassword = passwordgenerator.GenerateUserPassword("testpass-")
	user := &management.User{
		Username: testuser,
		Password: testpassword,
		Name:     testuser,
		Enabled:  &enabled,
	}

	newUser, err := users.CreateUserWithRole(client, user, "user")
	if err != nil {
		logrus.Fatalf("error creating admin client: %v", err)
	}

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	if err != nil {
		logrus.Fatalf("error creating standard user client: %v", err)
	}

	// create admin cluster
	adminClusterNames, err := createTestCluster(client, client, 1, "admintestcluster", cnis[0], kubernetesVersions[0], nodesAndRoles, advancedOptions)
	if err != nil {
		logrus.Fatalf("error creating admin user cluster: %v", err)
	}

	// create standard user clusters
	standardClusterNames, err := createTestCluster(standardUserClient, client, 2, "standardtestcluster", cnis[0], kubernetesVersions[0], nodesAndRoles, advancedOptions)
	if err != nil {
		logrus.Fatalf("error creating standard user clusters: %v", err)
	}

	//update userconfig
	userClusterConfig := UserClusterConfig{}
	userClusterConfig.Token = standardUserClient.Steve.Opts.TokenKey
	userClusterConfig.Username = newUser.Name
	userClusterConfig.AdminPassword = adminPassword
	userClusterConfig.Clusters = standardClusterNames

	rangerConfig.ClusterName = adminClusterNames[0]
	config.UpdateConfig(ranger.ConfigurationFileKey, rangerConfig)

	err = writeToConfigFile(userClusterConfig)
	if err != nil {
		logrus.Fatalf("error writing config file: %v", err)
	}
}

func getRangerKubeconfig() ([]byte, error) {
	kubeconfig, err := corral.GetKubeConfig(corralName)
	if err != nil {
		return nil, err
	}

	return kubeconfig, nil
}

func createAdminToken(password string, rangerConfig *ranger.Config) (string, error) {
	adminUser := &management.User{
		Username: "admin",
		Password: password,
	}

	hostURL := rangerConfig.Host
	var userToken *management.Token
	err := kwait.Poll(500*time.Millisecond, 5*time.Minute, func() (done bool, err error) {
		userToken, err = token.GenerateUserToken(adminUser, hostURL)
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return "", err
	}

	return userToken.Token, nil
}

func updateRangerDeployment(kubeconfig []byte) error {
	session := session.NewSession()
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	if err != nil {
		return err
	}

	restConfig, err := (clientConfig).ClientConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(session, restConfig)
	if err != nil {
		return err
	}

	deploymentResource := dynamicClient.Resource(deployments.DeploymentGroupVersionResource)

	cattleSystemDeploymentResource := deploymentResource.Namespace(namespace)
	unstructuredDeployment, err := cattleSystemDeploymentResource.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	updatedDeployment := &appv1.Deployment{}
	err = scheme.Scheme.Convert(unstructuredDeployment, updatedDeployment, unstructuredDeployment.GroupVersionKind())
	if err != nil {
		return err
	}

	updatedDeployment.Spec.Template.Spec.Containers[0].Args = []string{}
	updatedDeployment.Spec.Template.Spec.Containers[0].Image = rangerTestImage
	updatedDeployment.Spec.Strategy.RollingUpdate = nil
	updatedDeployment.Spec.Strategy.Type = appv1.RecreateDeploymentStrategyType

	_, err = cattleSystemDeploymentResource.Update(context.TODO(), unstructured.MustToUnstructured(updatedDeployment), metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	err = kwait.Poll(500*time.Millisecond, 10*time.Minute, func() (done bool, err error) {
		deployment, err := cattleSystemDeploymentResource.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if k8sErrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}

		newDeployment := &appv1.Deployment{}
		err = scheme.Scheme.Convert(deployment, newDeployment, deployment.GroupVersionKind())
		if err != nil {
			return false, err
		}
		if newDeployment.Status.ReadyReplicas == int32(3) {
			return true, nil
		}

		return false, nil
	})

	if err != nil {
		logrus.Infof("time out updating ranger deployment %v", err)
		return err
	}

	err = kwait.Poll(500*time.Millisecond, 10*time.Minute, func() (done bool, err error) {
		webhookDeployment, err := cattleSystemDeploymentResource.Get(context.TODO(), "ranger-webhook", metav1.GetOptions{})
		if k8sErrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}

		newDeployment := &appv1.Deployment{}
		err = scheme.Scheme.Convert(webhookDeployment, newDeployment, webhookDeployment.GroupVersionKind())
		if err != nil {
			return false, err
		}
		if newDeployment.Status.ReadyReplicas == int32(1) {
			return true, nil
		}
		return false, nil
	})
	logrus.Infof("time out updating ranger webhook deployment %v", err)
	return err
}

func createTestCluster(client, adminClient *ranger.Client, numClusters int, clusterNameBase, cni, kubeVersion string, nodesAndRoles []nodepools.NodeRoles, advancedOptions provisioning.AdvancedOptions) ([]string, error) {
	clusterNames := []string{}
	for i := 0; i < numClusters; i++ {
		nodeTemplateResp, err := aws.CreateAWSNodeTemplate(client)
		if err != nil {
			return nil, err
		}

		clusterName := namegen.AppendRandomString(clusterNameBase)
		clusterNames = append(clusterNames, clusterName)
		cluster := clusters.NewRKE1ClusterConfig(clusterName, cni, kubeVersion, "", client, advancedOptions)

		clusterResp, err := clusters.CreateRKE1Cluster(client, cluster)
		if err != nil {
			return nil, err
		}

		err = kwait.Poll(500*time.Millisecond, 10*time.Minute, func() (done bool, err error) {
			_, err = nodepools.NodePoolSetup(client, nodesAndRoles, clusterResp.ID, nodeTemplateResp.ID)
			if err != nil {
				return false, nil
			}
			return true, nil
		})
		if err != nil {
			return nil, err
		}

		opts := metav1.ListOptions{
			FieldSelector:  "metadata.name=" + clusterResp.ID,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		}
		watchInterface, err := adminClient.GetManagementWatchInterface(management.ClusterType, opts)
		if err != nil {
			return nil, err
		}

		checkFunc := clusters.IsHostedProvisioningClusterReady

		err = wait.WatchWait(watchInterface, checkFunc)
		if err != nil {
			return nil, err
		}
	}
	return clusterNames, nil
}

func writeToConfigFile(config UserClusterConfig) error {
	result := map[string]UserClusterConfig{}
	result[userClusterConfigsConfigurationFileKey] = config

	yamlConfig, err := yaml.Marshal(result)

	if err != nil {
		return err
	}

	return os.WriteFile("userclusterconfig.yaml", yamlConfig, 0644)
}
