package hostedcluster

import (
	"context"
	"fmt"
	"os"
	"strings"

	v3 "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/controllers/dashboard/chart"
	controllerv3 "github.com/ranger/ranger/pkg/generated/controllers/management.cattle.io/v3"
	controllerprojectv3 "github.com/ranger/ranger/pkg/generated/controllers/project.cattle.io/v3"
	"github.com/ranger/ranger/pkg/namespace"
	"github.com/ranger/ranger/pkg/project"
	"github.com/ranger/ranger/pkg/ref"
	"github.com/ranger/ranger/pkg/settings"
	"github.com/ranger/ranger/pkg/wrangler"
	v1 "github.com/ranger/wrangler/pkg/generated/controllers/core/v1"
	"github.com/ranger/wrangler/pkg/kv"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	AksCrdChart = chart.Definition{
		ReleaseNamespace: "cattle-system",
		ChartName:        "ranger-aks-operator-crd",
	}
	AksChart = chart.Definition{
		ReleaseNamespace: "cattle-system",
		ChartName:        "ranger-aks-operator",
	}
	EksCrdChart = chart.Definition{
		ReleaseNamespace: "cattle-system",
		ChartName:        "ranger-eks-operator-crd",
	}
	EksChart = chart.Definition{
		ReleaseNamespace: "cattle-system",
		ChartName:        "ranger-eks-operator",
	}
	GkeCrdChart = chart.Definition{
		ReleaseNamespace: "cattle-system",
		ChartName:        "ranger-gke-operator-crd",
	}
	GkeChart = chart.Definition{
		ReleaseNamespace: "cattle-system",
		ChartName:        "ranger-gke-operator",
	}
)

var (
	localCluster = "local"

	legacyOperatorAppNameFormat = "ranger-%s-operator"
)

type handler struct {
	manager      chart.Manager
	appCache     controllerprojectv3.AppCache
	apps         controllerprojectv3.AppController
	projectCache controllerv3.ProjectCache
	secretsCache v1.SecretCache
	chartsConfig chart.RangerConfigGetter
}

func Register(ctx context.Context, wContext *wrangler.Context) {
	h := &handler{
		manager:      wContext.SystemChartsManager,
		apps:         wContext.Project.App(),
		projectCache: wContext.Mgmt.Project().Cache(),
		secretsCache: wContext.Core.Secret().Cache(),
		appCache:     wContext.Project.App().Cache(),
		chartsConfig: chart.RangerConfigGetter{ConfigCache: wContext.Core.ConfigMap().Cache()},
	}

	wContext.Mgmt.Cluster().OnChange(ctx, "cluster-provisioning-operator", h.onClusterChange)
	wContext.Core.Secret().OnChange(ctx, "watch-helm-release", h.onSecretChange)
}

func (h handler) onClusterChange(key string, cluster *v3.Cluster) (*v3.Cluster, error) {
	if cluster == nil {
		return cluster, nil
	}

	skipChartInstallation := strings.EqualFold(settings.SkipHostedClusterChartInstallation.Get(), "true")
	if skipChartInstallation {
		logrus.Warn("Skipping installation of hosted cluster charts, 'skip-hosted-cluster-chart-installation' is set to true")
		return cluster, nil
	}

	var toInstallCrdChart, toInstallChart *chart.Definition
	var provider string
	if cluster.Spec.AKSConfig != nil {
		toInstallCrdChart = &AksCrdChart
		toInstallChart = &AksChart
		provider = "aks"
	} else if cluster.Spec.EKSConfig != nil {
		toInstallCrdChart = &EksCrdChart
		toInstallChart = &EksChart
		provider = "eks"
	} else if cluster.Spec.GKEConfig != nil {
		toInstallCrdChart = &GkeCrdChart
		toInstallChart = &GkeChart
		provider = "gke"
	}

	if toInstallCrdChart == nil || toInstallChart == nil {
		return cluster, nil
	}

	if err := h.removeLegacyOperatorIfExists(provider); err != nil {
		return cluster, err
	}

	if err := h.manager.Ensure(toInstallCrdChart.ReleaseNamespace, toInstallCrdChart.ChartName, "", "", nil, true, ""); err != nil {
		return cluster, err
	}

	systemGlobalRegistry := map[string]interface{}{
		"cattle": map[string]interface{}{
			"systemDefaultRegistry": settings.SystemDefaultRegistry.Get(),
		},
	}

	additionalCA, err := getAdditionalCA(h.secretsCache)
	if err != nil {
		return cluster, err
	}

	chartValues := map[string]interface{}{
		"global":               systemGlobalRegistry,
		"httpProxy":            os.Getenv("HTTP_PROXY"),
		"httpsProxy":           os.Getenv("HTTPS_PROXY"),
		"noProxy":              os.Getenv("NO_PROXY"),
		"additionalTrustedCAs": additionalCA != nil,
	}
	// add priority class value
	if priorityClassName, err := h.chartsConfig.GetPriorityClassName(); err != nil {
		if !apierror.IsNotFound(err) {
			logrus.Warnf("Failed to get ranger priorityClassName for %q: %v", toInstallChart.ChartName, err)
		}
	} else {
		chartValues[chart.PriorityClassKey] = priorityClassName
	}

	if err := h.manager.Ensure(toInstallChart.ReleaseNamespace, toInstallChart.ChartName, "", "", chartValues, true, ""); err != nil {
		return cluster, err
	}

	return cluster, nil
}

// check helm release secrets for aks/eks/gke operator chart, if it has been uninstalled, then remove it in m.manager.desiredChart
// so that we don't automatically redeploy it unless there is an AKS/EKS/GKE cluster triggering it
func (h handler) onSecretChange(key string, obj *corev1.Secret) (*corev1.Secret, error) {
	if obj == nil {
		ns, name := kv.Split(key, "/")
		if ns == namespace.System {
			// the name will follow the format sh.helm.release.v1.ranger-eks-operator-crd.v1
			parts := strings.Split(name, ".")
			if len(parts) == 6 {
				releaseName := parts[4]
				if isOperatorChartRelease(releaseName) {
					h.manager.Remove(ns, releaseName)
				}
			}
		}
	}
	return obj, nil
}

func isOperatorChartRelease(name string) bool {
	switch name {
	case AksCrdChart.ChartName, AksChart.ChartName, EksCrdChart.ChartName, EksChart.ChartName, GkeChart.ChartName, GkeCrdChart.ChartName:
		return true
	}
	return false
}

func (h handler) removeLegacyOperatorIfExists(provider string) error {
	systemProject, err := project.GetSystemProject(localCluster, h.projectCache)
	if err != nil {
		return err
	}

	systemProjectID := ref.Ref(systemProject)
	_, systemProjectName := ref.Parse(systemProjectID)

	legacyOperatorAppName := fmt.Sprintf(legacyOperatorAppNameFormat, provider)
	_, err = h.appCache.Get(systemProjectName, legacyOperatorAppName)
	if err != nil {
		if errors.IsNotFound(err) {
			// legacy app doesn't exist, no-op
			return nil
		}
		return err
	}

	return h.apps.Delete(systemProjectName, legacyOperatorAppName, &metav1.DeleteOptions{})
}

func getAdditionalCA(secretsCache v1.SecretCache) ([]byte, error) {
	secret, err := secretsCache.Get(namespace.System, "tls-ca-additional")
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	if secret == nil {
		return nil, nil
	}

	return secret.Data["ca-additional.pem"], nil
}
