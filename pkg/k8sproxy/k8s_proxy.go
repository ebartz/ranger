package k8sproxy

import (
	"net/http"

	"github.com/ranger/ranger/pkg/clusterrouter"
	"github.com/ranger/ranger/pkg/clusterrouter/proxy"
	"github.com/ranger/ranger/pkg/k8slookup"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/types/config/dialer"
)

func New(scaledContext *config.ScaledContext, dialer dialer.Factory, clusterContextGetter proxy.ClusterContextGetter) http.Handler {
	return clusterrouter.New(&scaledContext.RESTConfig, k8slookup.New(scaledContext, true), dialer,
		scaledContext.Management.Clusters("").Controller().Lister(),
		clusterContextGetter)
}
