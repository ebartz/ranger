package charts

import (
	"context"
	"fmt"

	"github.com/ranger/ranger/pkg/api/steve/catalog/types"
	catalogv1 "github.com/ranger/ranger/pkg/apis/catalog.cattle.io/v1"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	kubenamespaces "github.com/ranger/ranger/tests/framework/extensions/kubeapi/namespaces"
	"github.com/ranger/ranger/tests/framework/extensions/namespaces"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	// Namespace that ranger istio chart is installed in
	RangerIstioNamespace = "istio-system"
	// Name of the ranger istio chart
	RangerIstioName = "ranger-istio"
)

// InstallRangerIstioChart is a helper function that installs the ranger-istio chart.
func InstallRangerIstioChart(client *ranger.Client, installOptions *InstallOptions, rangerIstioOpts *RangerIstioOpts) error {
	serverSetting, err := client.Management.Setting.ByID(serverURLSettingID)
	if err != nil {
		return err
	}

	registrySetting, err := client.Management.Setting.ByID(defaultRegistrySettingID)
	if err != nil {
		return err
	}

	istioChartInstallActionPayload := &payloadOpts{
		InstallOptions:  *installOptions,
		Name:            RangerIstioName,
		Namespace:       RangerIstioNamespace,
		Host:            serverSetting.Value,
		DefaultRegistry: registrySetting.Value,
	}

	chartInstallAction := newIstioChartInstallAction(istioChartInstallActionPayload, rangerIstioOpts)

	catalogClient, err := client.GetClusterCatalogClient(installOptions.ClusterID)
	if err != nil {
		return err
	}

	// Cleanup registration
	client.Session.RegisterCleanupFunc(func() error {
		// UninstallAction for when uninstalling the ranger-istio chart
		defaultChartUninstallAction := newChartUninstallAction()

		err := catalogClient.UninstallChart(RangerIstioName, RangerIstioNamespace, defaultChartUninstallAction)
		if err != nil {
			return err
		}

		watchAppInterface, err := catalogClient.Apps(RangerIstioNamespace).Watch(context.TODO(), metav1.ListOptions{
			FieldSelector:  "metadata.name=" + RangerIstioName,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		})
		if err != nil {
			return err
		}

		err = wait.WatchWait(watchAppInterface, func(event watch.Event) (ready bool, err error) {
			if event.Type == watch.Error {
				return false, fmt.Errorf("there was an error uninstalling ranger istio chart")
			} else if event.Type == watch.Deleted {
				return true, nil
			}
			return false, nil
		})
		if err != nil {
			return err
		}

		steveclient, err := client.Steve.ProxyDownstream(installOptions.ClusterID)
		if err != nil {
			return err
		}

		namespaceClient := steveclient.SteveType(namespaces.NamespaceSteveType)

		namespace, err := namespaceClient.ByID(RangerIstioNamespace)
		if err != nil {
			return err
		}

		err = namespaceClient.Delete(namespace)
		if err != nil {
			return err
		}

		adminClient, err := ranger.NewClient(client.RangerConfig.AdminToken, client.Session)
		if err != nil {
			return err
		}
		adminDynamicClient, err := adminClient.GetDownStreamClusterClient(installOptions.ClusterID)
		if err != nil {
			return err
		}
		adminNamespaceResource := adminDynamicClient.Resource(kubenamespaces.NamespaceGroupVersionResource).Namespace("")

		watchNamespaceInterface, err := adminNamespaceResource.Watch(context.TODO(), metav1.ListOptions{
			FieldSelector:  "metadata.name=" + RangerIstioNamespace,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		})

		if err != nil {
			return err
		}

		return wait.WatchWait(watchNamespaceInterface, func(event watch.Event) (ready bool, err error) {
			if event.Type == watch.Deleted {
				return true, nil
			}
			return false, nil
		})
	})

	err = catalogClient.InstallChart(chartInstallAction)
	if err != nil {
		return err
	}

	// wait for chart to be full deployed
	watchAppInterface, err := catalogClient.Apps(RangerIstioNamespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + RangerIstioName,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	if err != nil {
		return err
	}

	err = wait.WatchWait(watchAppInterface, func(event watch.Event) (ready bool, err error) {
		app := event.Object.(*catalogv1.App)

		state := app.Status.Summary.State
		if state == string(catalogv1.StatusDeployed) {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

// newIstioChartInstallAction is a private helper function that returns chart install action with istio and payload options.
func newIstioChartInstallAction(p *payloadOpts, rangerIstioOpts *RangerIstioOpts) *types.ChartInstallAction {
	istioValues := map[string]interface{}{
		"tracing": map[string]interface{}{
			"enabled": rangerIstioOpts.Tracing,
		},
		"kiali": map[string]interface{}{
			"enabled": rangerIstioOpts.Kiali,
		},
		"ingressGateways": map[string]interface{}{
			"enabled": rangerIstioOpts.IngressGateways,
		},
		"egressGateways": map[string]interface{}{
			"enabled": rangerIstioOpts.EgressGateways,
		},
		"pilot": map[string]interface{}{
			"enabled": rangerIstioOpts.Pilot,
		},
		"telemetry": map[string]interface{}{
			"enabled": rangerIstioOpts.Telemetry,
		},
		"cni": map[string]interface{}{
			"enabled": rangerIstioOpts.CNI,
		},
	}
	chartInstall := newChartInstall(p.Name, p.InstallOptions.Version, p.InstallOptions.ClusterID, p.InstallOptions.ClusterName, p.Host, p.DefaultRegistry, istioValues)
	chartInstalls := []types.ChartInstall{*chartInstall}

	chartInstallAction := newChartInstallAction(p.Namespace, p.InstallOptions.ProjectID, chartInstalls)

	return chartInstallAction
}

// UpgradeRangerIstioChart is a helper function that upgrades the ranger-istio chart.
func UpgradeRangerIstioChart(client *ranger.Client, installOptions *InstallOptions, rangerIstioOpts *RangerIstioOpts) error {
	serverSetting, err := client.Management.Setting.ByID(serverURLSettingID)
	if err != nil {
		return err
	}

	registrySetting, err := client.Management.Setting.ByID(defaultRegistrySettingID)
	if err != nil {
		return err
	}

	istioChartUpgradeActionPayload := &payloadOpts{
		InstallOptions:  *installOptions,
		Name:            RangerIstioName,
		Namespace:       RangerIstioNamespace,
		Host:            serverSetting.Value,
		DefaultRegistry: registrySetting.Value,
	}

	chartUpgradeAction := newIstioChartUpgradeAction(istioChartUpgradeActionPayload, rangerIstioOpts)

	catalogClient, err := client.GetClusterCatalogClient(installOptions.ClusterID)
	if err != nil {
		return err
	}

	err = catalogClient.UpgradeChart(chartUpgradeAction)
	if err != nil {
		return err
	}

	adminClient, err := ranger.NewClient(client.RangerConfig.AdminToken, client.Session)
	if err != nil {
		return err
	}
	adminCatalogClient, err := adminClient.GetClusterCatalogClient(installOptions.ClusterID)
	if err != nil {
		return err
	}

	// wait for chart to be in status pending upgrade
	watchAppInterface, err := adminCatalogClient.Apps(RangerIstioNamespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + RangerIstioName,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	if err != nil {
		return err
	}

	err = wait.WatchWait(watchAppInterface, func(event watch.Event) (ready bool, err error) {
		app := event.Object.(*catalogv1.App)

		state := app.Status.Summary.State
		if state == string(catalogv1.StatusPendingUpgrade) {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	// wait for chart to be full deployed
	watchAppInterface, err = adminCatalogClient.Apps(RangerIstioNamespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + RangerIstioName,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	if err != nil {
		return err
	}

	err = wait.WatchWait(watchAppInterface, func(event watch.Event) (ready bool, err error) {
		app := event.Object.(*catalogv1.App)

		state := app.Status.Summary.State
		if state == string(catalogv1.StatusDeployed) {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}

// newIstioChartUpgradeAction is a private helper function that returns chart upgrade action with istio and payload options.
func newIstioChartUpgradeAction(p *payloadOpts, rangerIstioOpts *RangerIstioOpts) *types.ChartUpgradeAction {
	istioValues := map[string]interface{}{
		"tracing": map[string]interface{}{
			"enabled": rangerIstioOpts.Tracing,
		},
		"kiali": map[string]interface{}{
			"enabled": rangerIstioOpts.Kiali,
		},
		"ingressGateways": map[string]interface{}{
			"enabled": rangerIstioOpts.IngressGateways,
		},
		"egressGateways": map[string]interface{}{
			"enabled": rangerIstioOpts.EgressGateways,
		},
		"pilot": map[string]interface{}{
			"enabled": rangerIstioOpts.Pilot,
		},
		"telemetry": map[string]interface{}{
			"enabled": rangerIstioOpts.Telemetry,
		},
		"cni": map[string]interface{}{
			"enabled": rangerIstioOpts.CNI,
		},
	}
	chartUpgrade := newChartUpgrade(p.Name, p.InstallOptions.Version, p.InstallOptions.ClusterID, p.InstallOptions.ClusterName, p.Host, p.DefaultRegistry, istioValues)
	chartUpgrades := []types.ChartUpgrade{*chartUpgrade}

	chartUpgradeAction := newChartUpgradeAction(p.Namespace, chartUpgrades)

	return chartUpgradeAction
}
