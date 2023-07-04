package charts

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/ranger/ranger/pkg/controllers/managementuserlegacy/alert/config"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	v1 "github.com/ranger/ranger/tests/framework/clients/ranger/v1"
	"github.com/ranger/ranger/tests/framework/extensions/charts"
	"github.com/ranger/ranger/tests/framework/extensions/clusterrolebindings"
	"github.com/ranger/ranger/tests/framework/extensions/configmaps"
	"github.com/ranger/ranger/tests/framework/extensions/serviceaccounts"
	"github.com/ranger/ranger/tests/framework/extensions/workloads"
	"github.com/ranger/ranger/tests/framework/pkg/namegenerator"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	kubewait "k8s.io/apimachinery/pkg/util/wait"
)

const (
	// Project that example app and charts are installed in
	projectName = "System"
	// Secret path that contains encoded alert manager config
	secretPath = "alertmanager.yaml"
	// Secret ID that points alert manager secret
	alertManagerSecretId = charts.RangerMonitoringNamespace + "/" + charts.RangerMonitoringAlertSecret
	// Default random string length for random name generation
	defaultRandStringLength = 5
	// Webhook deployment annotation key that is being watched
	webhookReceiverAnnotationKey = "didReceiveRequestFromAlertmanager"
	// Webhook deployment annotation value that is being watched
	webhookReceiverAnnotationValue = "true"
	// Steve type for prometheus rules for schema
	prometheusRulesSteveType = "monitoring.coreos.com.prometheusrule"
	// rangerShellSettingID is the setting ID that used to grab ranger/shell image
	rangerShellSettingID = "shell-image"
	// Kubeconfig that linked to webhook deployment
	kubeConfig = `
apiVersion: v1
kind: Config
clusters:
- name: cluster
  cluster:
    certificate-authority: /run/secrets/kubernetes.io/serviceaccount/ca.crt
    server: https://kubernetes.default
contexts:
- name: default
  context:
    cluster: cluster
    user: user
current-context: default
users:
- name: user
  user:
    tokenFile: /run/secrets/kubernetes.io/serviceaccount/token
`
)

var (
	// Ranger monitoring chart alert manager path
	alertManagerPath = "api/v1/namespaces/cattle-monitoring-system/services/http:ranger-monitoring-alertmanager:9093/proxy/#/alerts"
	// Ranger monitoring chart grafana path
	grafanaPath = "api/v1/namespaces/cattle-monitoring-system/services/http:ranger-monitoring-grafana:80/proxy"
	// Ranger monitoring chart prometheus path
	prometheusPath = "api/v1/namespaces/cattle-monitoring-system/services/http:ranger-monitoring-prometheus:9090/proxy"
	// Ranger monitoring chart prometheus graph path
	prometheusGraphPath = prometheusPath + "/graph"
	// Ranger monitoring chart prometheus rules path
	prometheusRulesPath = prometheusPath + "/rules"
	// Ranger monitoring chart prometheus targets path
	prometheusTargetsPath = prometheusPath + "/targets"
	// Ranger monitoring chart prometheus targets API path
	prometheusTargetsPathAPI = prometheusPath + "/api/v1/targets"
	// Webhook receiver kubernetes object names
	webhookReceiverNamespaceName  = "webhook-namespace-" + namegenerator.RandStringLower(defaultRandStringLength)
	webhookReceiverDeploymentName = "webhook-" + namegenerator.RandStringLower(defaultRandStringLength)
	webhookReceiverServiceName    = "webhook-service-" + namegenerator.RandStringLower(defaultRandStringLength)
	// Label that is used to identify webhook and rule
	ruleLabel = map[string]string{"team": "qa"}
)

// waitUnknownPrometheusTargets is a private helper function
// that awaits the unknown Prometheus targets to be resolved until the timeout by using Prometheus API.
func waitUnknownPrometheusTargets(client *ranger.Client) error {
	checkUnknownPrometheusTargets := func() (bool, error) {
		var statusInit bool
		var unknownTargets []string
		resultAPI, err := charts.GetChartCaseEndpoint(client, client.RangerConfig.Host, prometheusTargetsPathAPI, true)
		if err != nil {
			return statusInit, err
		}
		var mapResponse map[string]interface{}
		if err = json.Unmarshal([]byte(resultAPI.Body), &mapResponse); err != nil {
			return statusInit, err
		}
		if mapResponse["status"] != "success" {
			return statusInit, errors.New("failed to get targets from prometheus")
		}
		activeTargets := mapResponse["data"].(map[string]interface{})["activeTargets"].([]interface{})
		if len(activeTargets) < 1 {
			return false, errors.New("failed to find any active targets")
		}
		for _, target := range activeTargets {
			targetMap := target.(map[string]interface{})
			if targetMap["health"].(string) == "unknown" {
				unknownTargets = append(unknownTargets, targetMap["labels"].(map[string]interface{})["instance"].(string))
			}
		}
		return len(unknownTargets) == 0, nil
	}

	err := kubewait.Poll(500*time.Millisecond, 2*time.Minute, func() (ongoing bool, err error) {
		result, err := checkUnknownPrometheusTargets()
		if err != nil {
			return ongoing, err
		}

		if result {
			return !ongoing, nil
		}

		return
	})
	if err != nil {
		return err
	}

	return nil
}

// checkPrometheusTargets is a private helper function
// that checks if all active prometheus targets are healthy by using prometheus API.
func checkPrometheusTargets(client *ranger.Client) (bool, error) {
	var statusInit bool
	var downTargets []string

	err := waitUnknownPrometheusTargets(client)
	if err != nil {
		return statusInit, err
	}

	resultAPI, err := charts.GetChartCaseEndpoint(client, client.RangerConfig.Host, prometheusTargetsPathAPI, true)
	if err != nil {
		return statusInit, err
	}

	var mapResponse map[string]interface{}
	if err = json.Unmarshal([]byte(resultAPI.Body), &mapResponse); err != nil {
		return statusInit, err
	}

	if mapResponse["status"] != "success" {
		return statusInit, errors.New("failed to get targets from prometheus")
	}

	activeTargets := mapResponse["data"].(map[string]interface{})["activeTargets"].([]interface{})
	if len(activeTargets) < 1 {
		return false, errors.New("failed to find any active targets")
	}

	for _, target := range activeTargets {
		targetMap := target.(map[string]interface{})
		if targetMap["health"].(string) == "down" {
			downTargets = append(downTargets, targetMap["labels"].(map[string]interface{})["instance"].(string))
		}
	}
	statusInit = len(downTargets) == 0

	if !statusInit {
		return statusInit, errors.Wrapf(err, "All active target(s) are not healthy: %v", downTargets)
	}

	return statusInit, nil
}

// editAlertReceiver is a private helper function
// that edits alert config structure to be used by the webhook receiver.
func editAlertReceiver(alertConfigByte []byte, origin string, originURL *url.URL) ([]byte, error) {
	alertConfig, err := config.Load(string(alertConfigByte))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal alert config")
	}

	alertConfig.Global = &config.GlobalConfig{
		ResolveTimeout: alertConfig.Global.ResolveTimeout,
	}
	alertConfig.Receivers = append(alertConfig.Receivers, &config.Receiver{
		Name: webhookReceiverDeploymentName,
		WebhookConfigs: []*config.WebhookConfig{
			{
				HTTPConfig: &config.HTTPClientConfig{
					ProxyURL: config.URL{URL: originURL},
				},
				NotifierConfig: config.NotifierConfig{
					VSendResolved: false,
				},
				URL: origin,
			},
		},
	})

	byteAlertConfig := []byte(alertConfig.String())
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(byteAlertConfig)))
	base64.StdEncoding.Encode(dst, byteAlertConfig)

	return byteAlertConfig, nil
}

// editAlertRoute is a private helper function
// that edits alert config structure to be used by the webhook receiver.
func editAlertRoute(alertConfigByte []byte, origin string, originURL *url.URL) ([]byte, error) {
	alertConfig, err := config.Load(string(alertConfigByte))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal alert config")
	}

	alertConfig.Global = &config.GlobalConfig{
		ResolveTimeout: alertConfig.Global.ResolveTimeout,
	}
	alertConfig.Route.Routes = append(alertConfig.Route.Routes, &config.Route{
		GroupWait:      alertConfig.Route.GroupWait,
		GroupInterval:  alertConfig.Route.GroupInterval,
		RepeatInterval: alertConfig.Route.RepeatInterval,
		Match:          ruleLabel,
		Receiver:       webhookReceiverDeploymentName,
	})

	byteAlertConfig := []byte(alertConfig.String())
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(byteAlertConfig)))
	base64.StdEncoding.Encode(dst, byteAlertConfig)

	return byteAlertConfig, nil
}

// createPrometheusRule is a private helper function
// that creates a prometheus rule to be used by the webhook receiver.
func createPrometheusRule(client *ranger.Client, clusterID string) error {
	ruleName := "webhook-rule-" + namegenerator.RandStringLower(defaultRandStringLength)
	alertName := "alert-" + namegenerator.RandStringLower(defaultRandStringLength)

	client.ReLogin()

	steveclient, err := client.Steve.ProxyDownstream(clusterID)
	if err != nil {
		return err
	}

	prometheusRule := &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ruleName,
			Namespace: charts.RangerMonitoringNamespace,
		},
		Spec: monitoringv1.PrometheusRuleSpec{
			Groups: []monitoringv1.RuleGroup{
				{
					Name: ruleName,
					Rules: []monitoringv1.Rule{
						{
							Alert:  alertName,
							Expr:   intstr.IntOrString{Type: intstr.String, StrVal: "vector(0)"},
							Labels: ruleLabel,
							For:    "0s",
						},
					},
				},
			},
		},
	}
	_, err = steveclient.SteveType(prometheusRulesSteveType).Create(prometheusRule)
	if err != nil {
		return err
	}

	return nil
}

// createWebhookReceiverDeployment is a private helper function that creates a service account, cluster role binding, and deployment for webhook receiver.
// The deployment has two different containers with a shared volume, one for kubectl commands, and the other one to receive requests and write access logs to the shared empty dir volume.
// Container that uses ranger/shell has a mounted volume to use the kubeconfig of the cluster. And it watches the access logs until a request from "alermanager" is received.
// When the request is received it sets its deployment annotation "didReceiveRequestFromAlertmanager" to "true" while the annotations being watched by the test itself.
func createAlertWebhookReceiverDeployment(client *ranger.Client, clusterID, namespace, deploymentName string) (*v1.SteveAPIObject, error) {
	serviceAccountName := "alert-receiver-sa-" + namegenerator.RandStringLower(defaultRandStringLength)
	clusterRoleBindingName := "alert-receiver-cluster-admin-" + namegenerator.RandStringLower(defaultRandStringLength)
	configMapName := "alert-receiver-cm-" + namegenerator.RandStringLower(defaultRandStringLength)

	steveclient, err := client.Steve.ProxyDownstream(clusterID)
	if err != nil {
		return nil, err
	}

	// Create webhook receiver service account
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace,
		},
	}

	_, err = steveclient.SteveType(serviceaccounts.ServiceAccountSteveType).Create(serviceAccount)
	if err != nil {
		return nil, err
	}

	// Create webhook receiver cluster role binding
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccount.Name,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}

	_, err = steveclient.SteveType(clusterrolebindings.ClusterRoleBindingSteveType).Create(clusterRoleBinding)
	if err != nil {
		return nil, err
	}

	// Create webhook receiver config map
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"config": kubeConfig,
		},
	}

	_, err = steveclient.SteveType(configmaps.ConfigMapSteveType).Create(configmap)
	if err != nil {
		return nil, err
	}

	labels := map[string]string{}
	labels["workload.user.cattle.io/workloadselector"] = fmt.Sprintf("apps.deployment-%v-%v", namespace, deploymentName)

	imageSetting, err := client.Management.Setting.ByID(rangerShellSettingID)
	if err != nil {
		return nil, err
	}

	// Create webhook receiver deployment
	var runAsUser int64
	var runAsGroup int64
	podSpecTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "alert-reciver-deployment",
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:      corev1.RestartPolicyNever,
			ServiceAccountName: serviceAccount.Name,
			Containers: []corev1.Container{
				{
					Name:    "kubectl",
					Image:   imageSetting.Value,
					Command: []string{"/bin/sh", "-c"},
					Args: []string{
						fmt.Sprintf(
							`until [ "$didReceiveRequestFromAlertmanager" = true ]; do if grep -q "Alertmanager" "/traefik/access.log"; then kubectl patch deployment %s -n %s --type "json" -p '[{"op":"add","path":"/metadata/annotations/%s","value":"%s"}]'; didReceiveRequestFromAlertmanager=true; sleep 5m; else sleep 10; echo "Checking logs file one more time"; fi; done`,
							deploymentName, namespace, webhookReceiverAnnotationKey, webhookReceiverAnnotationValue,
						),
					},
					SecurityContext: &corev1.SecurityContext{
						RunAsUser:  &runAsUser,
						RunAsGroup: &runAsGroup,
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "config", MountPath: "/root/usr/share/.kube/"},
						{Name: "logs", MountPath: "/traefik"},
					},
				},
				{
					Name:  "traefik",
					Image: "traefik:latest",
					Args: []string{
						"--entrypoints.web.address=:80", "--api.dashboard=true", "--api.insecure=true", "--accesslog=true", "--accesslog.filepath=/var/log/traefik/access.log", "--log.level=INFO", "--accesslog.fields.headers.defaultmode=keep",
					},
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
							Protocol:      corev1.ProtocolTCP,
						},
						{
							ContainerPort: 8080,
							Protocol:      corev1.ProtocolTCP,
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "logs", MountPath: "/var/log/traefik"},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: configmap.Name},
						},
					},
				},
				{
					Name: "logs",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	isCattleLabeled := true
	deploymentTemplate := workloads.NewDeploymentTemplate(deploymentName, namespace, podSpecTemplate, isCattleLabeled, nil)
	deployment, err := steveclient.SteveType(workloads.DeploymentSteveType).Create(deploymentTemplate)
	if err != nil {
		return deployment, err
	}

	return deployment, nil
}
