package provisioningv2

import (
	"context"

	"github.com/ranger/ranger/pkg/controllers/provisioningv2/cluster"
	"github.com/ranger/ranger/pkg/controllers/provisioningv2/fleetcluster"
	"github.com/ranger/ranger/pkg/controllers/provisioningv2/fleetworkspace"
	"github.com/ranger/ranger/pkg/controllers/provisioningv2/managedchart"
	"github.com/ranger/ranger/pkg/controllers/provisioningv2/provisioningcluster"
	"github.com/ranger/ranger/pkg/controllers/provisioningv2/provisioninglog"
	"github.com/ranger/ranger/pkg/controllers/provisioningv2/secret"
	"github.com/ranger/ranger/pkg/features"
	"github.com/ranger/ranger/pkg/provisioningv2/kubeconfig"
	"github.com/ranger/ranger/pkg/wrangler"
)

func Register(ctx context.Context, clients *wrangler.Context, kubeconfigManager *kubeconfig.Manager) {
	cluster.Register(ctx, clients, kubeconfigManager)
	if features.MCM.Enabled() {
		secret.Register(ctx, clients)
	}
	provisioningcluster.Register(ctx, clients)
	provisioninglog.Register(ctx, clients)

	if features.Fleet.Enabled() {
		managedchart.Register(ctx, clients)
		fleetcluster.Register(ctx, clients)
		fleetworkspace.Register(ctx, clients)
	}
}
