package customresourcedefinitions

import (
	"context"

	"github.com/ranger/ranger/tests/framework/clients/ranger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// helper function that returns a pointer to an unstructured list of custom resource definitions
func ListCustomResourceDefinitions(client *ranger.Client, clusterID string, namespace string) (*unstructured.UnstructuredList, error) {
	dynamicClient, err := client.GetDownStreamClusterClient(clusterID)
	if err != nil {
		return nil, err
	}

	customResourceDefinitionResource := dynamicClient.Resource(CustomResourceDefinitions).Namespace(namespace)
	CRDs, err := customResourceDefinitionResource.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return CRDs, err
}
