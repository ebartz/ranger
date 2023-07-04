package secrets

import (
	"context"

	"github.com/ranger/ranger/pkg/api/scheme"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/unstructured"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
)

// CreateSecretForCluster is a helper function that uses the ranger client to create a secret in a namespace for a specific cluster.
func CreateSecretForCluster(client *ranger.Client, secret *corev1.Secret, clusterID, namespace string) (*corev1.Secret, error) {
	dynamicClient, err := client.GetDownStreamClusterClient(clusterID)
	if err != nil {
		return nil, err
	}
	secretResource := dynamicClient.Resource(SecretGroupVersionResource).Namespace(namespace)

	return CreateSecret(secretResource, secret)
}

// CreateSecret is a helper function that uses the dynamic client to create a secret in a namespace for a specific cluster.
func CreateSecret(secretResource dynamic.ResourceInterface, secret *corev1.Secret) (*corev1.Secret, error) {
	unstructuredResp, err := secretResource.Create(context.TODO(), unstructured.MustToUnstructured(secret), metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	newSecret := &corev1.Secret{}
	err = scheme.Scheme.Convert(unstructuredResp, newSecret, unstructuredResp.GroupVersionKind())
	if err != nil {
		return nil, err
	}
	return newSecret, nil
}
