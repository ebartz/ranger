package management

import (
	"context"

	"github.com/ranger/ranger/pkg/clustermanager"
	"github.com/ranger/ranger/pkg/controllers/management/aks"
	"github.com/ranger/ranger/pkg/controllers/management/authprovisioningv2"
	"github.com/ranger/ranger/pkg/controllers/management/clusterupstreamrefresher"
	"github.com/ranger/ranger/pkg/controllers/management/eks"
	"github.com/ranger/ranger/pkg/controllers/management/feature"
	"github.com/ranger/ranger/pkg/controllers/management/gke"
	"github.com/ranger/ranger/pkg/controllers/management/k3sbasedupgrade"
	"github.com/ranger/ranger/pkg/features"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/wrangler"
)

func RegisterWrangler(ctx context.Context, wranglerContext *wrangler.Context, management *config.ManagementContext, manager *clustermanager.Manager) error {
	k3sbasedupgrade.Register(ctx, wranglerContext, management, manager)
	aks.Register(ctx, wranglerContext, management)
	eks.Register(ctx, wranglerContext, management)
	gke.Register(ctx, wranglerContext, management)
	clusterupstreamrefresher.Register(ctx, wranglerContext)

	feature.Register(ctx, wranglerContext)

	if features.ProvisioningV2.Enabled() {
		if err := authprovisioningv2.Register(ctx, wranglerContext, management); err != nil {
			return err
		}
	}

	return nil
}
