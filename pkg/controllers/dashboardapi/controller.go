package dashboardapi

import (
	"context"

	"github.com/ranger/ranger/pkg/controllers/dashboard/helm"
	"github.com/ranger/ranger/pkg/controllers/dashboardapi/feature"
	"github.com/ranger/ranger/pkg/controllers/dashboardapi/settings"
	"github.com/ranger/ranger/pkg/wrangler"
)

func Register(ctx context.Context, wrangler *wrangler.Context) error {
	feature.Register(ctx, wrangler.Mgmt.Feature())
	helm.RegisterReposForFollowers(ctx, wrangler.Core.Secret().Cache(), wrangler.Catalog.ClusterRepo())
	return settings.Register(wrangler.Mgmt.Setting())
}
