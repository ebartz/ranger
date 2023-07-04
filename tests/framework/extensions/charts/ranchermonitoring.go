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
	// Namespace that ranger monitoring chart is installed in
	RangerMonitoringNamespace = "cattle-monitoring-system"
	// Name of the ranger monitoring chart
	RangerMonitoringName = "ranger-monitoring"
	// Name of the ranger monitoring alert config secret
	RangerMonitoringAlertSecret = "alertmanager-ranger-monitoring-alertmanager"
	// Name of ranger monitoring crd chart
	RangerMonitoringCRDName = "ranger-monitoring-crd"
)

// InstallRangerMonitoringChart is a helper function that installs the ranger-monitoring chart.
func InstallRangerMonitoringChart(client *ranger.Client, installOptions *InstallOptions, rangerMonitoringOpts *RangerMonitoringOpts) error {
	serverSetting, err := client.Management.Setting.ByID(serverURLSettingID)
	if err != nil {
		return err
	}

	registrySetting, err := client.Management.Setting.ByID(defaultRegistrySettingID)
	if err != nil {
		return err
	}

	monitoringChartInstallActionPayload := &payloadOpts{
		InstallOptions:  *installOptions,
		Name:            RangerMonitoringName,
		Namespace:       RangerMonitoringNamespace,
		Host:            serverSetting.Value,
		DefaultRegistry: registrySetting.Value,
	}

	chartInstallAction := newMonitoringChartInstallAction(monitoringChartInstallActionPayload, rangerMonitoringOpts)

	catalogClient, err := client.GetClusterCatalogClient(installOptions.ClusterID)
	if err != nil {
		return err
	}

	// Cleanup registration
	client.Session.RegisterCleanupFunc(func() error {
		// UninstallAction for when uninstalling the ranger-monitoring chart
		defaultChartUninstallAction := newChartUninstallAction()

		err = catalogClient.UninstallChart(RangerMonitoringName, RangerMonitoringNamespace, defaultChartUninstallAction)
		if err != nil {
			return err
		}

		watchAppInterface, err := catalogClient.Apps(RangerMonitoringNamespace).Watch(context.TODO(), metav1.ListOptions{
			FieldSelector:  "metadata.name=" + RangerMonitoringName,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		})
		if err != nil {
			return err
		}

		err = wait.WatchWait(watchAppInterface, func(event watch.Event) (ready bool, err error) {
			if event.Type == watch.Error {
				return false, fmt.Errorf("there was an error uninstalling ranger monitoring chart")
			} else if event.Type == watch.Deleted {
				return true, nil
			}
			return false, nil
		})
		if err != nil {
			return err
		}

		err = catalogClient.UninstallChart(RangerMonitoringCRDName, RangerMonitoringNamespace, defaultChartUninstallAction)
		if err != nil {
			return err
		}

		watchAppInterface, err = catalogClient.Apps(RangerMonitoringNamespace).Watch(context.TODO(), metav1.ListOptions{
			FieldSelector:  "metadata.name=" + RangerMonitoringCRDName,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		})
		if err != nil {
			return err
		}

		err = wait.WatchWait(watchAppInterface, func(event watch.Event) (ready bool, err error) {
			chart := event.Object.(*catalogv1.App)
			if event.Type == watch.Error {
				return false, fmt.Errorf("there was an error uninstalling ranger monitoring chart")
			} else if event.Type == watch.Deleted {
				return true, nil
			} else if chart == nil {
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

		namespace, err := namespaceClient.ByID(RangerMonitoringNamespace)
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
			FieldSelector:  "metadata.name=" + RangerMonitoringNamespace,
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
	watchAppInterface, err := catalogClient.Apps(RangerMonitoringNamespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + RangerMonitoringName,
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

// newMonitoringChartInstallAction is a private helper function that returns chart install action with monitoring and payload options.
func newMonitoringChartInstallAction(p *payloadOpts, rangerMonitoringOpts *RangerMonitoringOpts) *types.ChartInstallAction {
	monitoringValues := map[string]interface{}{
		"ingressNgnix": map[string]interface{}{
			"enabled": rangerMonitoringOpts.IngressNginx,
		},
		"prometheus": map[string]interface{}{
			"prometheusSpec": map[string]interface{}{
				"evaluationInterval": "1m",
				"retentionSize":      "50GiB",
				"scrapeInterval":     "1m",
			},
		},
		"rkeControllerManager": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEControllerManager,
		},
		"rkeEtcd": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEEtcd,
		},
		"rkeProxy": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEProxy,
		},
		"rkeScheduler": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEScheduler,
		},
	}

	chartInstall := newChartInstall(p.Name, p.InstallOptions.Version, p.InstallOptions.ClusterID, p.InstallOptions.ClusterName, p.Host, p.DefaultRegistry, monitoringValues)
	chartInstallCRD := newChartInstall(p.Name+"-crd", p.Version, p.InstallOptions.ClusterID, p.InstallOptions.ClusterName, p.Host, p.DefaultRegistry, nil)
	chartInstalls := []types.ChartInstall{*chartInstallCRD, *chartInstall}

	chartInstallAction := newChartInstallAction(p.Namespace, p.ProjectID, chartInstalls)

	return chartInstallAction
}

// UpgradeMonitoringChart is a helper function that upgrades the ranger-monitoring chart.
func UpgradeRangerMonitoringChart(client *ranger.Client, installOptions *InstallOptions, rangerMonitoringOpts *RangerMonitoringOpts) error {
	serverSetting, err := client.Management.Setting.ByID(serverURLSettingID)
	if err != nil {
		return err
	}

	registrySetting, err := client.Management.Setting.ByID(defaultRegistrySettingID)
	if err != nil {
		return err
	}

	monitoringChartUpgradeActionPayload := &payloadOpts{
		InstallOptions:  *installOptions,
		Name:            RangerMonitoringName,
		Namespace:       RangerMonitoringNamespace,
		Host:            serverSetting.Value,
		DefaultRegistry: registrySetting.Value,
	}

	chartUpgradeAction := newMonitoringChartUpgradeAction(monitoringChartUpgradeActionPayload, rangerMonitoringOpts)

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
	watchAppInterface, err := adminCatalogClient.Apps(RangerMonitoringNamespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + RangerMonitoringName,
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
	watchAppInterface, err = adminCatalogClient.Apps(RangerMonitoringNamespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + RangerMonitoringName,
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

// newMonitoringChartUpgradeAction is a private helper function that returns chart upgrade action with monitoring and payload options.
func newMonitoringChartUpgradeAction(p *payloadOpts, rangerMonitoringOpts *RangerMonitoringOpts) *types.ChartUpgradeAction {
	monitoringValues := map[string]interface{}{
		"ingressNgnix": map[string]interface{}{
			"enabled": rangerMonitoringOpts.IngressNginx,
		},
		"prometheus": map[string]interface{}{
			"prometheusSpec": map[string]interface{}{
				"evaluationInterval": "1m",
				"retentionSize":      "50GiB",
				"scrapeInterval":     "1m",
			},
		},
		"rkeControllerManager": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEControllerManager,
		},
		"rkeEtcd": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEEtcd,
		},
		"rkeProxy": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEProxy,
		},
		"rkeScheduler": map[string]interface{}{
			"enabled": rangerMonitoringOpts.RKEScheduler,
		},
	}
	chartUpgrade := newChartUpgrade(p.Name, p.InstallOptions.Version, p.InstallOptions.ClusterID, p.InstallOptions.ClusterName, p.Host, p.DefaultRegistry, monitoringValues)
	chartUpgradeCRD := newChartUpgrade(p.Name+"-crd", p.InstallOptions.Version, p.InstallOptions.ClusterID, p.InstallOptions.ClusterName, p.Host, p.DefaultRegistry, monitoringValues)
	chartUpgrades := []types.ChartUpgrade{*chartUpgradeCRD, *chartUpgrade}

	chartUpgradeAction := newChartUpgradeAction(p.Namespace, chartUpgrades)

	return chartUpgradeAction
}
