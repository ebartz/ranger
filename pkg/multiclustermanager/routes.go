package multiclustermanager

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ranger/apiserver/pkg/parse"
	"github.com/ranger/ranger/pkg/api/norman"
	"github.com/ranger/ranger/pkg/api/norman/customization/aks"
	"github.com/ranger/ranger/pkg/api/norman/customization/clusterregistrationtokens"
	"github.com/ranger/ranger/pkg/api/norman/customization/gke"
	"github.com/ranger/ranger/pkg/api/norman/customization/oci"
	"github.com/ranger/ranger/pkg/api/norman/customization/vsphere"
	managementapi "github.com/ranger/ranger/pkg/api/norman/server"
	"github.com/ranger/ranger/pkg/api/steve/supportconfigs"
	"github.com/ranger/ranger/pkg/auth/providers/publicapi"
	"github.com/ranger/ranger/pkg/auth/providers/saml"
	"github.com/ranger/ranger/pkg/auth/requests"
	"github.com/ranger/ranger/pkg/auth/requests/sar"
	"github.com/ranger/ranger/pkg/auth/tokens"
	"github.com/ranger/ranger/pkg/auth/webhook"
	"github.com/ranger/ranger/pkg/channelserver"
	"github.com/ranger/ranger/pkg/clustermanager"
	rangerdialer "github.com/ranger/ranger/pkg/dialer"
	"github.com/ranger/ranger/pkg/httpproxy"
	k8sProxyPkg "github.com/ranger/ranger/pkg/k8sproxy"
	"github.com/ranger/ranger/pkg/metrics"
	"github.com/ranger/ranger/pkg/multiclustermanager/whitelist"
	"github.com/ranger/ranger/pkg/rbac"
	"github.com/ranger/ranger/pkg/rkenodeconfigserver"
	"github.com/ranger/ranger/pkg/telemetry"
	"github.com/ranger/ranger/pkg/tunnelserver/mcmauthorizer"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/version"
	"github.com/ranger/steve/pkg/auth"
)

func router(ctx context.Context, localClusterEnabled bool, tunnelAuthorizer *mcmauthorizer.Authorizer, scaledContext *config.ScaledContext, clusterManager *clustermanager.Manager) (func(http.Handler) http.Handler, error) {
	var (
		k8sProxy             = k8sProxyPkg.New(scaledContext, scaledContext.Dialer, clusterManager)
		connectHandler       = scaledContext.Dialer.(*rangerdialer.Factory).TunnelServer
		connectConfigHandler = rkenodeconfigserver.Handler(tunnelAuthorizer, scaledContext)
		clusterImport        = clusterregistrationtokens.ClusterImport{Clusters: scaledContext.Management.Clusters("")}
	)

	tokenAPI, err := tokens.NewAPIHandler(ctx, scaledContext, norman.ConfigureAPIUI)
	if err != nil {
		return nil, err
	}

	publicAPI, err := publicapi.NewHandler(ctx, scaledContext, norman.ConfigureAPIUI)
	if err != nil {
		return nil, err
	}

	managementAPI, err := managementapi.New(ctx, scaledContext, clusterManager, k8sProxy, localClusterEnabled)
	if err != nil {
		return nil, err
	}

	metaProxy, err := httpproxy.NewProxy("/proxy/", whitelist.Proxy.Get, scaledContext)
	if err != nil {
		return nil, err
	}

	metricsHandler := metrics.NewMetricsHandler(scaledContext, clusterManager, promhttp.Handler())

	channelserver := channelserver.NewHandler(ctx)

	supportConfigGenerator := supportconfigs.NewHandler(scaledContext)
	// Unauthenticated routes
	unauthed := mux.NewRouter()
	unauthed.UseEncodedPath()

	unauthed.Path("/").MatcherFunc(parse.MatchNotBrowser).Handler(managementAPI)
	unauthed.Handle("/v3/connect/config", connectConfigHandler)
	unauthed.Handle("/v3/connect", connectHandler)
	unauthed.Handle("/v3/connect/register", connectHandler)
	unauthed.Handle("/v3/import/{token}_{clusterId}.yaml", http.HandlerFunc(clusterImport.ClusterImportHandler))
	unauthed.Handle("/v3/settings/cacerts", managementAPI).MatcherFunc(onlyGet)
	unauthed.Handle("/v3/settings/first-login", managementAPI).MatcherFunc(onlyGet)
	unauthed.Handle("/v3/settings/ui-banners", managementAPI).MatcherFunc(onlyGet)
	unauthed.Handle("/v3/settings/ui-issues", managementAPI).MatcherFunc(onlyGet)
	unauthed.Handle("/v3/settings/ui-pl", managementAPI).MatcherFunc(onlyGet)
	unauthed.Handle("/v3/settings/ui-brand", managementAPI).MatcherFunc(onlyGet)
	unauthed.Handle("/v3/settings/ui-default-landing", managementAPI).MatcherFunc(onlyGet)
	unauthed.Handle("/rangerversion", version.NewVersionHandler())
	unauthed.PathPrefix("/v1-{prefix}-release/channel").Handler(channelserver)
	unauthed.PathPrefix("/v1-{prefix}-release/release").Handler(channelserver)
	unauthed.PathPrefix("/v1-saml").Handler(saml.AuthHandler())
	unauthed.PathPrefix("/v3-public").Handler(publicAPI)

	// Authenticated routes
	authed := mux.NewRouter()
	authed.UseEncodedPath()
	impersonatingAuth := auth.ToMiddleware(requests.NewImpersonatingAuth(sar.NewSubjectAccessReview(clusterManager)))
	accessControlHandler := rbac.NewAccessControlHandler()

	authed.Use(mux.MiddlewareFunc(impersonatingAuth))
	authed.Use(mux.MiddlewareFunc(accessControlHandler))
	authed.Use(requests.NewAuthenticatedFilter)

	authed.Path("/meta/{resource:aks.+}").Handler(aks.NewAKSHandler(scaledContext))
	authed.Path("/meta/{resource:gke.+}").Handler(gke.NewGKEHandler(scaledContext))
	authed.Path("/meta/oci/{resource}").Handler(oci.NewOCIHandler(scaledContext))
	authed.Path("/meta/vsphere/{field}").Handler(vsphere.NewVsphereHandler(scaledContext))
	authed.Path("/v3/tokenreview").Methods(http.MethodPost).Handler(&webhook.TokenReviewer{})
	authed.Path("/metrics/{clusterID}").Handler(metricsHandler)
	authed.Path(supportconfigs.Endpoint).Handler(&supportConfigGenerator)
	authed.PathPrefix("/k8s/clusters/").Handler(k8sProxy)
	authed.PathPrefix("/meta/proxy").Handler(metaProxy)
	authed.PathPrefix("/v1-telemetry").Handler(telemetry.NewProxy())
	authed.PathPrefix("/v3/identit").Handler(tokenAPI)
	authed.PathPrefix("/v3/token").Handler(tokenAPI)
	authed.PathPrefix("/v3").Handler(managementAPI)

	// Metrics authenticated route
	metricsAuthed := mux.NewRouter()
	metricsAuthed.UseEncodedPath()
	tokenReviewAuth := auth.ToMiddleware(requests.NewTokenReviewAuth(scaledContext.K8sClient.AuthenticationV1()))
	metricsAuthed.Use(mux.MiddlewareFunc(tokenReviewAuth.Chain(impersonatingAuth)))
	metricsAuthed.Use(mux.MiddlewareFunc(accessControlHandler))
	metricsAuthed.Use(requests.NewAuthenticatedFilter)

	metricsAuthed.Path("/metrics").Handler(metricsHandler)

	unauthed.NotFoundHandler = authed
	authed.NotFoundHandler = metricsAuthed
	return func(next http.Handler) http.Handler {
		metricsAuthed.NotFoundHandler = next
		return unauthed
	}, nil
}

// onlyGet will match only GET but will not return a 405 like route.Methods and instead just not match
func onlyGet(req *http.Request, m *mux.RouteMatch) bool {
	return req.Method == http.MethodGet
}
