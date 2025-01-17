package deployments

import (
	"context"
	"fmt"

	"github.com/ranger/ranger/pkg/api/scheme"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/unstructured"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DeploymentGroupVersionResource is the required Group Version Resource for accessing deployments in a cluster,
// using the dynamic client.
var DeploymentGroupVersionResource = schema.GroupVersionResource{
	Group:    "apps",
	Version:  "v1",
	Resource: "deployments",
}

// CreateDeployment is a helper function that uses the dynamic client to create a deployment on a namespace for a specific cluster.
func CreateDeployment(client *ranger.Client, clusterName, deploymentName, namespace string, template corev1.PodTemplateSpec) (*appv1.Deployment, error) {
	dynamicClient, err := client.GetDownStreamClusterClient(clusterName)
	if err != nil {
		return nil, err
	}

	labels := map[string]string{}
	labels["workload.user.cattle.io/workloadselector"] = fmt.Sprintf("apps.deployment-%v-%v", namespace, deploymentName)

	template.ObjectMeta = metav1.ObjectMeta{
		Labels: labels,
	}

	template.Spec.RestartPolicy = corev1.RestartPolicyAlways
	deployment := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: namespace,
		},
		Spec: appv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: template,
		},
	}

	deploymentResource := dynamicClient.Resource(DeploymentGroupVersionResource).Namespace(namespace)

	unstructuredResp, err := deploymentResource.Create(context.TODO(), unstructured.MustToUnstructured(deployment), metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	newDeployment := &appv1.Deployment{}
	err = scheme.Scheme.Convert(unstructuredResp, newDeployment, unstructuredResp.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	return newDeployment, nil
}
