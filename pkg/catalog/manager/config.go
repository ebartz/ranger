package manager

import (
	v3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
)

type CatalogInfo struct {
	catalog        *v3.Catalog
	projectCatalog *v3.ProjectCatalog
	clusterCatalog *v3.ClusterCatalog
}
