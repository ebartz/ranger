package primechecks

import (
	"fmt"
	"strings"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	client "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/ranger/ranger/tests/framework/extensions/rangerversion"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const (
	PodResourceSteveType = "pod"
	rangerImage         = "ranger"
)

// CheckUIBrand checks the UI brand of Ranger Prime. If the Ranger instance is not Ranger Prime, the UI brand should be blank.
func CheckUIBrand(client *ranger.Client, isPrime bool, rangerBrand *client.Setting, brand string) error {
	if isPrime && brand != rangerBrand.Value {
		return fmt.Errorf("error: Ranger Prime UI brand %s does not match defined UI brand %s", rangerBrand.Value, brand)
	}

	return nil
}

// CheckVersion checks the if Ranger Prime is set to true and the version of Ranger.
func CheckVersion(isPrime bool, rangerVersion string, serverConfig *rangerversion.Config) error {
	if isPrime && rangerVersion != serverConfig.RangerVersion {
		return fmt.Errorf("error: Ranger Prime: %t | Version: %s", isPrime, serverConfig.RangerVersion)
	}

	return nil
}

// CheckSystemDefaultRegistry checks if the system default registry is set to the expected value.
func CheckSystemDefaultRegistry(isPrime bool, primeRegistry string, registry *client.Setting) error {
	if isPrime && primeRegistry != registry.Value {
		return fmt.Errorf("error: Ranger Prime system default registry %s does not match user defined registry %s", registry.Value, primeRegistry)
	}

	return nil
}

// CheckLocalClusterRangerImages checks if the Ranger images are set to the expected registry.
func CheckLocalClusterRangerImages(client *ranger.Client, isPrime bool, rangerVersion, primeRegistry, clusterID string) ([]string, []error) {
	downstreamClient, err := client.Steve.ProxyDownstream(clusterID)
	if err != nil {
		return nil, []error{err}
	}

	steveClient := downstreamClient.SteveType(PodResourceSteveType)

	pods, err := steveClient.List(nil)
	if err != nil {
		return nil, []error{err}
	}

	var imageResults []string
	var imageErrors []error

	for _, pod := range pods.Data {
		podStatus := &corev1.PodStatus{}
		err = v1.ConvertToK8sType(pod.Status, podStatus)
		if err != nil {
			return nil, []error{err}
		}

		image := podStatus.ContainerStatuses[0].Image

		if (strings.Contains(image, primeRegistry) && isPrime) || (strings.Contains(image, rangerImage) && !isPrime) {
			imageResults = append(imageResults, fmt.Sprintf("INFO: %s: %s\n", pod.Name, image))
			logrus.Infof("Pod %s is using image: %s", pod.Name, image)
		} else if strings.Contains(image, rangerImage) && isPrime {
			imageErrors = append(imageErrors, fmt.Errorf("ERROR: %s: %s", pod.Name, image))
			logrus.Infof("Pod %s is using image: %s", pod.Name, image)
		}
	}

	return imageResults, imageErrors
}
