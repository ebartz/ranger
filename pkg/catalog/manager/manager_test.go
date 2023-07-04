package manager

import (
	"testing"

	"github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/settings"
	"github.com/stretchr/testify/assert"
)

func TestLatestAvailableTemplateVersion(t *testing.T) {
	template := &v3.CatalogTemplate{
		Spec: v3.TemplateSpec{
			Versions: []v3.TemplateVersionSpec{
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.16",
					Version:           "0.12.16",
					RangerMinVersion: "v2.2.0",
					RangerMaxVersion: "v2.3.0",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.15",
					Version:           "0.12.15",
					RangerMinVersion: "v2.1.0",
					RangerMaxVersion: "v2.2.0",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.14",
					Version:           "0.12.14",
					RangerMinVersion: "v2.0.0",
					RangerMaxVersion: "v2.1.0",
				},
			},
		},
	}

	templateWithoutRangerVersion := &v3.CatalogTemplate{
		Spec: v3.TemplateSpec{
			Versions: []v3.TemplateVersionSpec{
				{
					ExternalID: "catalog://?catalog=library&template=artifactory-ha&version=0.12.16",
					Version:    "0.12.16",
				},
				{
					ExternalID: "catalog://?catalog=library&template=artifactory-ha&version=0.12.15",
					Version:    "0.12.15",
				},
				{
					ExternalID: "catalog://?catalog=library&template=artifactory-ha&version=0.12.14",
					Version:    "0.12.14",
				},
			},
		},
	}

	templateWithoutMinRangerVersion := &v3.CatalogTemplate{
		Spec: v3.TemplateSpec{
			Versions: []v3.TemplateVersionSpec{
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.16",
					Version:           "0.12.16",
					RangerMaxVersion: "v2.3.0",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.15",
					Version:           "0.12.15",
					RangerMaxVersion: "v2.2.0",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.14",
					Version:           "0.12.14",
					RangerMaxVersion: "v2.1.0",
				},
			},
		},
	}

	templateWithoutMaxRangerVersion := &v3.CatalogTemplate{
		Spec: v3.TemplateSpec{
			Versions: []v3.TemplateVersionSpec{
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.16",
					Version:           "0.12.16",
					RangerMinVersion: "v2.2.0",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.15",
					Version:           "0.12.15",
					RangerMinVersion: "v2.1.0",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.14",
					Version:           "0.12.14",
					RangerMinVersion: "v2.0.0",
				},
			},
		},
	}

	templateWithMinPrerelease := &v3.CatalogTemplate{
		Spec: v3.TemplateSpec{
			Versions: []v3.TemplateVersionSpec{
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.16",
					Version:           "0.12.16",
					RangerMinVersion: "v2.2.1-rc1",
					RangerMaxVersion: "v2.2.2",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.15",
					Version:           "0.12.15",
					RangerMinVersion: "v2.2.0-rc1",
					RangerMaxVersion: "v2.2.1",
				},
				{
					ExternalID: "catalog://?catalog=library&template=artifactory-ha&version=0.12.14",
					Version:    "0.12.14",
				},
			},
		},
	}

	templateWith99MaxPatchVersion := &v3.CatalogTemplate{
		Spec: v3.TemplateSpec{
			Versions: []v3.TemplateVersionSpec{
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.16",
					Version:           "0.12.16",
					RangerMinVersion: "v2.3.0-alpha1",
					RangerMaxVersion: "v2.3.99",
				},
				{
					ExternalID:        "catalog://?catalog=library&template=artifactory-ha&version=0.12.15",
					Version:           "0.12.15",
					RangerMinVersion: "v2.2.0-alpha1",
					RangerMaxVersion: "v2.2.99",
				},
				{
					ExternalID: "catalog://?catalog=library&template=artifactory-ha&version=0.12.14",
					Version:    "0.12.14",
				},
			},
		},
	}

	testLatestAvailableTemplateVersion(t, "v2.0.1", "0.12.14", template)
	testLatestAvailableTemplateVersion(t, "v2.0.2-beta", "0.12.14", template)
	testLatestAvailableTemplateVersion(t, "v2.1.0-alpha1", "0.12.14", template)
	testLatestAvailableTemplateVersion(t, "v2.1.0", "0.12.15", template)
	testLatestAvailableTemplateVersion(t, "v2.2.5", "0.12.16", template)
	testLatestAvailableTemplateVersion(t, "dev", "0.12.16", template)
	testLatestAvailableTemplateVersion(t, "master", "0.12.16", template)
	testLatestAvailableTemplateVersion(t, "master-head", "0.12.16", template)
	testLatestAvailableTemplateVersion(t, "", "0.12.16", template)

	testLatestAvailableTemplateVersion(t, "v2.0.1", "0.12.16", templateWithoutRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.0.2-beta", "0.12.16", templateWithoutRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.1.0-alpha1", "0.12.16", templateWithoutRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.1.0", "0.12.16", templateWithoutRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.2.5", "0.12.16", templateWithoutRangerVersion)
	testLatestAvailableTemplateVersion(t, "master", "0.12.16", templateWithoutRangerVersion)
	testLatestAvailableTemplateVersion(t, "master-head", "0.12.16", templateWithoutRangerVersion)
	testLatestAvailableTemplateVersion(t, "", "0.12.16", templateWithoutRangerVersion)

	testLatestAvailableTemplateVersion(t, "v2.0.1", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.0.2-beta", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.1.0-alpha1", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.1.0", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.2.5", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "dev", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "master", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "master-head", "0.12.16", templateWithoutMinRangerVersion)
	testLatestAvailableTemplateVersion(t, "", "0.12.16", templateWithoutMinRangerVersion)

	testLatestAvailableTemplateVersion(t, "v2.0.1", "0.12.14", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.0.2-beta", "0.12.14", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.1.0-alpha1", "0.12.14", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.1.0", "0.12.15", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "v2.2.5", "0.12.16", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "dev", "0.12.16", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "master", "0.12.16", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "master-head", "0.12.16", templateWithoutMaxRangerVersion)
	testLatestAvailableTemplateVersion(t, "", "0.12.16", templateWithoutMaxRangerVersion)

	testLatestAvailableTemplateVersion(t, "v2.2.0-0", "0.12.14", templateWithMinPrerelease)
	testLatestAvailableTemplateVersion(t, "v2.2.0-alpha1", "0.12.14", templateWithMinPrerelease)
	testLatestAvailableTemplateVersion(t, "v2.2.0-rc1", "0.12.15", templateWithMinPrerelease)
	testLatestAvailableTemplateVersion(t, "v2.2.0", "0.12.15", templateWithMinPrerelease)
	testLatestAvailableTemplateVersion(t, "v2.2.1-0", "0.12.15", templateWithMinPrerelease)
	testLatestAvailableTemplateVersion(t, "v2.2.1-alpha1", "0.12.15", templateWithMinPrerelease)
	testLatestAvailableTemplateVersion(t, "v2.2.1-rc1", "0.12.16", templateWithMinPrerelease)
	testLatestAvailableTemplateVersion(t, "v2.2.1", "0.12.16", templateWithMinPrerelease)

	testLatestAvailableTemplateVersion(t, "v2.2.0-0", "0.12.14", templateWith99MaxPatchVersion)
	testLatestAvailableTemplateVersion(t, "v2.2.0-alpha1", "0.12.15", templateWith99MaxPatchVersion)
	testLatestAvailableTemplateVersion(t, "v2.2.0-rc1", "0.12.15", templateWith99MaxPatchVersion)
	testLatestAvailableTemplateVersion(t, "v2.2.0", "0.12.15", templateWith99MaxPatchVersion)
	testLatestAvailableTemplateVersion(t, "v2.3.0-0", "0.12.14", templateWith99MaxPatchVersion)
	testLatestAvailableTemplateVersion(t, "v2.3.0-alpha1", "0.12.16", templateWith99MaxPatchVersion)
	testLatestAvailableTemplateVersion(t, "v2.3.0-rc1", "0.12.16", templateWith99MaxPatchVersion)
	testLatestAvailableTemplateVersion(t, "v2.3.0", "0.12.16", templateWith99MaxPatchVersion)
}

func testLatestAvailableTemplateVersion(t *testing.T, serverVersion, expectedCatalogVersion string, template *v3.CatalogTemplate) {
	err := settings.ServerVersion.Set(serverVersion)
	assert.Nil(t, err)

	catalogManager := Manager{}
	templateVersion, err := catalogManager.LatestAvailableTemplateVersion(template, "")
	assert.Nil(t, err)
	assert.Equal(t, expectedCatalogVersion, templateVersion.Version)
}
