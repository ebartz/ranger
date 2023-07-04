package dashboard

import (
	"context"

	"github.com/ranger/ranger/pkg/controllers/capi"
	"github.com/ranger/ranger/pkg/controllers/capr"
	"github.com/ranger/ranger/pkg/controllers/dashboard/apiservice"
	"github.com/ranger/ranger/pkg/controllers/dashboard/clusterindex"
	"github.com/ranger/ranger/pkg/controllers/dashboard/clusterregistrationtoken"
	"github.com/ranger/ranger/pkg/controllers/dashboard/cspadaptercharts"
	"github.com/ranger/ranger/pkg/controllers/dashboard/fleetcharts"
	"github.com/ranger/ranger/pkg/controllers/dashboard/helm"
	"github.com/ranger/ranger/pkg/controllers/dashboard/hostedcluster"
	"github.com/ranger/ranger/pkg/controllers/dashboard/kubernetesprovider"
	"github.com/ranger/ranger/pkg/controllers/dashboard/mcmagent"
	"github.com/ranger/ranger/pkg/controllers/dashboard/scaleavailable"
	"github.com/ranger/ranger/pkg/controllers/dashboard/systemcharts"
	"github.com/ranger/ranger/pkg/controllers/management/clusterconnected"
	"github.com/ranger/ranger/pkg/controllers/provisioningv2"
	"github.com/ranger/ranger/pkg/features"
	"github.com/ranger/ranger/pkg/provisioningv2/kubeconfig"
	"github.com/ranger/ranger/pkg/wrangler"
	"github.com/ranger/wrangler/pkg/needacert"
	"github.com/sirupsen/logrus"
)

func Register(ctx context.Context, wrangler *wrangler.Context, embedded bool, registryOverride string) error {
	helm.Register(ctx, wrangler)
	kubernetesprovider.Register(ctx,
		wrangler.Mgmt.Cluster(),
		wrangler.K8s,
		wrangler.MultiClusterManager)
	apiservice.Register(ctx, wrangler, embedded)
	needacert.Register(ctx,
		wrangler.Core.Secret(),
		wrangler.Core.Service(),
		wrangler.Admission.MutatingWebhookConfiguration(),
		wrangler.Admission.ValidatingWebhookConfiguration(),
		wrangler.CRD.CustomResourceDefinition())
	scaleavailable.Register(ctx, wrangler)
	if err := systemcharts.Register(ctx, wrangler, registryOverride); err != nil {
		return err
	}

	if err := cspadaptercharts.Register(ctx, wrangler); err != nil {
		return err
	}

	clusterconnected.Register(ctx, wrangler)

	if features.MCM.Enabled() {
		hostedcluster.Register(ctx, wrangler)
	}

	if features.Fleet.Enabled() {
		if err := fleetcharts.Register(ctx, wrangler); err != nil {
			return err
		}
	}

	if features.ProvisioningV2.Enabled() || features.MCM.Enabled() {
		clusterregistrationtoken.Register(ctx, wrangler)
	}

	if features.ProvisioningV2.Enabled() {
		kubeconfigManager := kubeconfig.New(wrangler)
		clusterindex.Register(ctx, wrangler)
		provisioningv2.Register(ctx, wrangler, kubeconfigManager)
		if features.RKE2.Enabled() {
			capr.Register(ctx, wrangler, kubeconfigManager)
		}
	}

	if features.EmbeddedClusterAPI.Enabled() {
		capiStart, err := capi.Register(ctx, wrangler)
		if err != nil {
			return err
		}
		wrangler.OnLeader(func(ctx context.Context) error {
			if err := capiStart(ctx); err != nil {
				logrus.Fatal(err)
			}
			logrus.Info("Cluster API is started")
			return nil
		})
	}

	if features.MCMAgent.Enabled() || features.MCM.Enabled() {
		err := mcmagent.Register(ctx, wrangler)
		if err != nil {
			return err
		}
	}

	return nil
}
