package ingresses

import (
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/charts"
)

const (
	IngressSteveType = "networking.k8s.io.ingress"
)

// AccessIngressExternally checks if the ingress is accessible externally,
// it returns true if the ingress is accessible, false if it is not, and an error if there is an error.
func AccessIngressExternally(client *ranger.Client, hostname string, isWithTLS bool) (bool, error) {
	result, err := charts.GetChartCaseEndpoint(client, hostname, "", isWithTLS)
	if err != nil {
		return false, err
	}

	return result.Ok, err
}
