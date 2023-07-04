package catalog

import (
	"fmt"
	"strings"
	"testing"

	apimgmtv3 "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/catalog/manager"
	mgmtv3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3/fakes"
	"github.com/ranger/ranger/pkg/settings"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/version"
)

type testCase struct {
	externalID        string
	rangerVersion    string
	kubernetesVersion string
	result            bool
}

func newTestCase(templateName, rangerVersion, kubernetesVersion string, result bool) testCase {
	return testCase{
		externalID:        generateExternalID(templateName),
		rangerVersion:    rangerVersion,
		kubernetesVersion: kubernetesVersion,
		result:            result,
	}
}
func TestExtractLinks(t *testing.T) {
	templateStore := setupTestTemplateStore()

	// Setup test cases
	var tests []testCase
	tests = append(tests, newTestCase("ranger-testRangerConstraint-0.1", "2.4.0", "", true))
	tests = append(tests, newTestCase("ranger-testRangerConstraint-0.1", "2.5.1", "", false))
	// all test cases without ranger version should result in true, since a dev environment is assumed and checks are bypassed
	tests = append(tests, newTestCase("ranger-testRangerConstraint-0.1", "", "", true))
	tests = append(tests, newTestCase("ranger-testKubeVersionConstraint-0.1", "2.4.0", "1.14.0", true))
	tests = append(tests, newTestCase("ranger-testKubeVersionConstraint-0.1", "2.4.0", "1.16.0", false))
	tests = append(tests, newTestCase("ranger-testKubeVersionConstraint-0.1", "2.4.0", "", true))
	tests = append(tests, newTestCase("ranger-testBothConstraints-0.1", "2.5.1", "1.14.0", false))
	tests = append(tests, newTestCase("ranger-testBothConstraints-0.1", "2.4.9", "1.16.0", false))
	tests = append(tests, newTestCase("ranger-testBothConstraints-0.1", "2.5.1", "1.16.0", false))
	tests = append(tests, newTestCase("ranger-testBothConstraints-0.1", "2.4.9", "1.14.0", true))
	tests = append(tests, newTestCase("ranger-testBothConstraints-0.1", "", "", true))
	tests = append(tests, newTestCase("ranger-testBothConstraints-0.1", "", "1.16.0", true))
	tests = append(tests, newTestCase("ranger-testNeitherConstraints-0.1", "", "", true))
	tests = append(tests, newTestCase("ranger-testNeitherConstraints-0.1", "2.5.1", "1.16.0", true))

	// The externalID of a testCase will be tested for compatibility with the given rangerVersion and cluster kubernetes
	// version. The result field defines whether compatibility is expected.
	for _, test := range tests {
		settings.ServerVersion.Set(test.rangerVersion)
		templateStore.CatalogManager = &manager.Manager{ClusterLister: newClusterListerWithVersion(test.kubernetesVersion)}
		clusterName := "test"
		if test.kubernetesVersion == "" {
			// kubernetesVersion being empty is equivalent to not targeting a cluster
			clusterName = ""
		}
		assert.Equal(t, test.result, templateStore.isTemplateVersionCompatible(clusterName, test.externalID), fmt.Sprintf("kubeVersion [%s] or rangerVersion [%s] unexpected compatibility with template [%s]", test.kubernetesVersion, test.rangerVersion, test.externalID))
	}
}

// setupTestTemplateStore sets up test templates with different combinations of ranger
// version constraints and kubeVersionConstraints
func setupTestTemplateStore() templateStore {
	catalogTemplateVersionLister := &fakes.CatalogTemplateVersionListerMock{
		GetFunc: func(namespace, name string) (*mgmtv3.CatalogTemplateVersion, error) {
			switch name {
			case "ranger-testRangerConstraint-0.1":
				return newCatalogTemplateVersion("2.4.99", ""), nil
			case "ranger-testKubeVersionConstraint-0.1":
				return newCatalogTemplateVersion("", "<1.15.0"), nil
			case "ranger-testBothConstraints-0.1":
				return newCatalogTemplateVersion("2.4.99", "<1.15.0"), nil
			case "ranger-testNeitherConstraints-0.1":
				return newCatalogTemplateVersion("", ""), nil
			}
			return nil, nil
		},
	}

	return templateStore{
		CatalogTemplateVersionLister: catalogTemplateVersionLister,
		CatalogManager:               &manager.Manager{},
	}
}
func newClusterListerWithVersion(kubernetesVersion string) *fakes.ClusterListerMock {
	return &fakes.ClusterListerMock{
		GetFunc: func(namespace, name string) (*mgmtv3.Cluster, error) {
			if name == "test" {
				cluster := mgmtv3.Cluster{
					Status: apimgmtv3.ClusterStatus{
						Version: &version.Info{
							GitVersion: kubernetesVersion,
						},
					},
				}
				return &cluster, nil
			}
			return nil, fmt.Errorf("invalid cluster: %s", name)
		},
	}
}

func newCatalogTemplateVersion(maxRangerVersion, kubeVersion string) *apimgmtv3.CatalogTemplateVersion {
	catalogTemplateVersion := &apimgmtv3.CatalogTemplateVersion{
		Spec: apimgmtv3.TemplateVersionSpec{
			KubeVersion:       kubeVersion,
			RangerMaxVersion: maxRangerVersion,
		},
	}
	return catalogTemplateVersion
}

func generateExternalID(name string) string {
	parts := strings.Split(name, "-")
	return fmt.Sprintf("catalog://?catalog=%s&template=%s&version=%s", parts[0], parts[1], parts[2])
}
