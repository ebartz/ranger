package codecoverage

import (
	"context"
	"fmt"
	"strings"
	"time"

	apiv1 "github.com/ranger/ranger/pkg/apis/provisioning.cattle.io/v1"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/ranger/ranger/tests/framework/extensions/clusters"
	"github.com/ranger/ranger/tests/framework/extensions/kubeconfig"
	"github.com/ranger/ranger/tests/framework/pkg/killserver"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
)

var podGroupVersionResource = corev1.SchemeGroupVersion.WithResource("pods")

const (
	cattleSystemNameSpace = "cattle-system"
	localCluster          = "local"
	rangerCoverFile      = "rangercoverage"
	agentCoverFile        = "agentcoverage"
	outputDir             = "cover"
)

func checkServiceIsRunning(dynamicClient dynamic.Interface) error {
	return kwait.Poll(500*time.Millisecond, 2*time.Minute, func() (done bool, err error) {
		_, err = dynamicClient.Resource(podGroupVersionResource).Namespace(cattleSystemNameSpace).List(context.Background(), metav1.ListOptions{})
		if k8sErrors.IsInternalError(err) || k8sErrors.IsServiceUnavailable(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}
		return true, nil
	})
}

func killTestServices(client *ranger.Client, clusterID string, podNames []string) error {
	cmd := []string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("curl -s localhost%s", killserver.KillServerPort),
	}

	kubeConfig, err := kubeconfig.GetKubeconfig(client, clusterID)
	if err != nil {
		return err
	}

	restConfig, err := (*kubeConfig).ClientConfig()
	if err != nil {
		return err
	}

	for _, podName := range podNames {
		_, err := kubeconfig.KubectlExec(restConfig, podName, cattleSystemNameSpace, cmd)
		if err != nil {
			logrus.Errorf("error killing pod container %v", err)
		}
	}

	return nil
}

func retrieveCodeCoverageFile(client *ranger.Client, clusterID, coverageFilename string, podNames []string) error {
	kubeConfig, err := kubeconfig.GetKubeconfig(client, clusterID)
	if err != nil {
		return err
	}

	restConfig, err := (*kubeConfig).ClientConfig()
	if err != nil {
		return err
	}

	for _, podName := range podNames {
		fileName := fmt.Sprintf("%s%s", podName, coverageFilename)
		dst := fmt.Sprintf("%s/%s", outputDir, fileName)

		err := kubeconfig.CopyFileFromPod(restConfig, *kubeConfig, podName, cattleSystemNameSpace, coverageFilename, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

// KillRangerTestServicesRetrieveCoverage is a function that kills the ranger service of the local cluster
// inorder for the code coverage report to be written, and then copies over the coverage reports from the pods
// to a local destination. The custom code coverage ranger image must be running in the local cluster.
func KillRangerTestServicesRetrieveCoverage(client *ranger.Client) error {
	var podNames []string
	dynamicClient, err := client.GetRangerDynamicClient()
	if err != nil {
		return err
	}

	pods, err := dynamicClient.Resource(podGroupVersionResource).Namespace(cattleSystemNameSpace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		name := pod.GetName()
		if strings.Contains(name, "ranger") && !strings.Contains(name, "webhook") {
			podNames = append(podNames, pod.GetName())
		}
	}

	err = killTestServices(client, localCluster, podNames)
	if err != nil {
		return err
	}

	err = checkServiceIsRunning(dynamicClient)
	if err != nil {
		return err
	}

	return retrieveCodeCoverageFile(client, localCluster, rangerCoverFile, podNames)
}

// KillAgentTestServicesRetrieveCoverage is a function that kills the cattle-cluster-agent service of a downstream cluster
// inorder for the code coverage report to be written, and then copies over the coverage reports from the pods
// to a local destination. The custom code coverage ranger-agent image must be running in the downstream cluster.
func KillAgentTestServicesRetrieveCoverage(client *ranger.Client) error {
	clusters, err := client.Steve.SteveType(clusters.ProvisioningSteveResourceType).ListAll(nil)
	if err != nil {
		return err
	}

	for _, cluster := range clusters.Data {
		clusterStatus := &apiv1.ClusterStatus{}
		err = v1.ConvertToK8sType(cluster.Status, clusterStatus)
		if err != nil {
			return err
		}
		clusterID := clusterStatus.ClusterName
		if clusterID != localCluster {
			dynamicClient, err := client.GetDownStreamClusterClient(clusterID)
			if err != nil {
				logrus.Errorf("could not connect to downstream cluster")
				continue
			}

			pods, err := dynamicClient.Resource(podGroupVersionResource).Namespace(cattleSystemNameSpace).List(context.Background(), metav1.ListOptions{})
			if err != nil {
				logrus.Errorf("could not list pods")
				continue
			}

			var podNames []string
			for _, pod := range pods.Items {
				if strings.Contains(pod.GetName(), "cattle-cluster-agent") {
					podNames = append(podNames, pod.GetName())
				}
			}

			err = killTestServices(client, clusterID, podNames)
			if err != nil {
				return err
			}

			err = checkServiceIsRunning(dynamicClient)
			if err != nil {
				return err
			}

			err = retrieveCodeCoverageFile(client, clusterID, agentCoverFile, podNames)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
