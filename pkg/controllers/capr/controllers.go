package capr

import (
	"context"

	"github.com/ranger/ranger/pkg/capr"
	"github.com/ranger/ranger/pkg/capr/planner"
	"github.com/ranger/ranger/pkg/controllers/capr/bootstrap"
	"github.com/ranger/ranger/pkg/controllers/capr/dynamicschema"
	"github.com/ranger/ranger/pkg/controllers/capr/machinedrain"
	"github.com/ranger/ranger/pkg/controllers/capr/machinenodelookup"
	"github.com/ranger/ranger/pkg/controllers/capr/machineprovision"
	"github.com/ranger/ranger/pkg/controllers/capr/managesystemagent"
	plannercontroller "github.com/ranger/ranger/pkg/controllers/capr/planner"
	"github.com/ranger/ranger/pkg/controllers/capr/plansecret"
	"github.com/ranger/ranger/pkg/controllers/capr/rkecluster"
	"github.com/ranger/ranger/pkg/controllers/capr/rkecontrolplane"
	"github.com/ranger/ranger/pkg/controllers/capr/unmanaged"
	"github.com/ranger/ranger/pkg/features"
	"github.com/ranger/ranger/pkg/provisioningv2/image"
	"github.com/ranger/ranger/pkg/provisioningv2/kubeconfig"
	"github.com/ranger/ranger/pkg/provisioningv2/systeminfo"
	"github.com/ranger/ranger/pkg/settings"
	"github.com/ranger/ranger/pkg/wrangler"
)

func Register(ctx context.Context, clients *wrangler.Context, kubeconfigManager *kubeconfig.Manager) {
	rkePlanner := planner.New(ctx, clients, planner.InfoFunctions{
		ImageResolver:           image.ResolveWithControlPlane,
		ReleaseData:             capr.GetKDMReleaseData,
		SystemAgentImage:        settings.SystemAgentInstallerImage.Get,
		SystemPodLabelSelectors: systeminfo.NewRetriever(clients).GetSystemPodLabelSelectors,
	})
	if features.MCM.Enabled() {
		dynamicschema.Register(ctx, clients)
		machineprovision.Register(ctx, clients, kubeconfigManager)
	}
	rkecluster.Register(ctx, clients)
	bootstrap.Register(ctx, clients)
	machinenodelookup.Register(ctx, clients, kubeconfigManager)
	plannercontroller.Register(ctx, clients, rkePlanner)
	plansecret.Register(ctx, clients)
	unmanaged.Register(ctx, clients, kubeconfigManager)
	rkecontrolplane.Register(ctx, clients)
	managesystemagent.Register(ctx, clients)
	machinedrain.Register(ctx, clients)
}
