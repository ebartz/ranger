package managementagent

import (
	"context"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/ranger/norman/store/crd"
	"github.com/ranger/norman/types"
	projectclient "github.com/ranger/ranger/pkg/client/generated/project/v3"
	"github.com/ranger/ranger/pkg/controllers/managementagent/dnsrecord"
	"github.com/ranger/ranger/pkg/controllers/managementagent/endpoints"
	"github.com/ranger/ranger/pkg/controllers/managementagent/externalservice"
	"github.com/ranger/ranger/pkg/controllers/managementagent/ingress"
	"github.com/ranger/ranger/pkg/controllers/managementagent/ingresshostgen"
	"github.com/ranger/ranger/pkg/controllers/managementagent/nslabels"
	"github.com/ranger/ranger/pkg/controllers/managementagent/nsserviceaccount"
	"github.com/ranger/ranger/pkg/controllers/managementagent/podresources"
	"github.com/ranger/ranger/pkg/controllers/managementagent/servicemonitor"
	"github.com/ranger/ranger/pkg/controllers/managementagent/targetworkloadservice"
	"github.com/ranger/ranger/pkg/controllers/managementagent/workload"
	"github.com/ranger/ranger/pkg/controllers/managementuserlegacy/monitoring"
	"github.com/ranger/ranger/pkg/features"
	pkgmonitoring "github.com/ranger/ranger/pkg/monitoring"
	"github.com/ranger/ranger/pkg/schemas/factory"
	"github.com/ranger/ranger/pkg/types/config"
)

func Register(ctx context.Context, cluster *config.UserOnlyContext) error {
	dnsrecord.Register(ctx, cluster)
	externalservice.Register(ctx, cluster)
	endpoints.Register(ctx, cluster)
	ingress.Register(ctx, cluster)
	ingresshostgen.Register(ctx, cluster)
	nslabels.Register(ctx, cluster)
	podresources.Register(ctx, cluster)
	targetworkloadservice.Register(ctx, cluster)
	workload.Register(ctx, cluster)
	nsserviceaccount.Register(ctx, cluster)

	if features.MonitoringV1.Enabled() {
		if err := createUserClusterCRDs(ctx, cluster); err != nil {
			return err
		}
		servicemonitor.Register(ctx, cluster)
		monitoring.RegisterAgent(ctx, cluster)
	}

	return nil
}

func createUserClusterCRDs(ctx context.Context, c *config.UserOnlyContext) error {
	overrided := struct {
		types.Namespaced
	}{}

	schemas := factory.Schemas(&pkgmonitoring.APIVersion).
		MustImport(&pkgmonitoring.APIVersion, monitoringv1.Prometheus{}, overrided).
		MustImport(&pkgmonitoring.APIVersion, monitoringv1.PrometheusRule{}, overrided).
		MustImport(&pkgmonitoring.APIVersion, monitoringv1.ServiceMonitor{}, overrided).
		MustImport(&pkgmonitoring.APIVersion, monitoringv1.Alertmanager{}, overrided)

	f, err := crd.NewFactoryFromClient(c.RESTConfig)
	if err != nil {
		return err
	}

	_, err = f.CreateCRDs(ctx, config.UserStorageContext,
		schemas.Schema(&pkgmonitoring.APIVersion, projectclient.PrometheusType),
		schemas.Schema(&pkgmonitoring.APIVersion, projectclient.PrometheusRuleType),
		schemas.Schema(&pkgmonitoring.APIVersion, projectclient.AlertmanagerType),
		schemas.Schema(&pkgmonitoring.APIVersion, projectclient.ServiceMonitorType),
	)

	f.BatchWait()

	return err
}
