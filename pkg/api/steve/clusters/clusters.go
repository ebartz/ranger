package clusters

import (
	"context"
	"net/http"
	"time"

	"github.com/ranger/apiserver/pkg/handlers"
	"github.com/ranger/apiserver/pkg/types"
	"github.com/ranger/ranger/pkg/api/steve/norman"
	"github.com/ranger/ranger/pkg/auth/providers/common"
	"github.com/ranger/ranger/pkg/auth/requests"
	"github.com/ranger/ranger/pkg/clusterrouter"
	normanv3 "github.com/ranger/ranger/pkg/schemas/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/settings"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/wrangler"
	"github.com/ranger/steve/pkg/podimpersonation"
	schema2 "github.com/ranger/steve/pkg/schema"
	steve "github.com/ranger/steve/pkg/server"
	"github.com/ranger/wrangler/pkg/schemas"
	"github.com/ranger/wrangler/pkg/schemas/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Register(ctx context.Context, server *steve.Server, wrangler *wrangler.Context) error {
	log := &log{
		cg: server.ClientFactory,
	}
	shell := &shell{
		cg:              server.ClientFactory,
		namespace:       "cattle-system",
		impersonator:    podimpersonation.New("shell", server.ClientFactory, time.Hour, settings.FullShellImage),
		clusterRegistry: server.ClusterRegistry,
	}
	sc, err := config.NewScaledContext(*wrangler.RESTConfig, nil)
	if err != nil {
		return err
	}

	userManager, err := common.NewUserManagerNoBindings(sc)
	if err != nil {
		return err
	}
	kubeconfig := kubeconfigDownload{
		userMgr: userManager,
		auth:    requests.NewAuthenticator(ctx, clusterrouter.GetClusterID, sc),
	}

	server.ClusterCache.OnAdd(ctx, shell.impersonator.PurgeOldRoles)
	server.ClusterCache.OnChange(ctx, func(gvk schema.GroupVersionKind, key string, obj, oldObj runtime.Object) error {
		return shell.impersonator.PurgeOldRoles(gvk, key, obj)
	})

	server.BaseSchemas.MustImportAndCustomize(GenerateKubeconfigOutput{}, nil)
	server.SchemaFactory.AddTemplate(schema2.Template{
		Group:     "management.cattle.io",
		Kind:      "Cluster",
		Formatter: norman.NewLinksAndActionsFormatter(wrangler.MultiClusterManager, normanv3.Version, "cluster"),
		Customize: func(schema *types.APISchema) {
			if schema.LinkHandlers == nil {
				schema.LinkHandlers = map[string]http.Handler{}
			}
			schema.LinkHandlers["shell"] = shell
			schema.LinkHandlers["log"] = log
			if schema.ActionHandlers == nil {
				schema.ActionHandlers = map[string]http.Handler{}
			}
			schema.ActionHandlers["generateKubeconfig"] = kubeconfig
			if schema.ResourceActions == nil {
				schema.ResourceActions = map[string]schemas.Action{}
			}
			schema.ResourceActions["generateKubeconfig"] = schemas.Action{
				Output: "generateKubeconfigOutput",
			}
			schema.ByIDHandler = func(request *types.APIRequest) (types.APIObject, error) {
				// By pass authorization for local shell because the user might not have
				// GET granted for local cluster
				if request.Name == "local" && request.Link == "shell" {
					shell.ServeHTTP(request.Response, request.Request)
					return types.APIObject{}, validation.ErrComplete
				}
				return handlers.ByIDHandler(request)
			}
			// Everybody can list even if they have no list or get privileges. The users
			// authorization will still be used to determine what can be seen but just
			// may result in an empty list
			schema.CollectionMethods = append(schema.CollectionMethods, http.MethodGet)
		},
	})
	server.SchemaFactory.AddTemplate(schema2.Template{
		Group: "management.cattle.io",
		Kind:  "Project",
		Customize: func(schema *types.APISchema) {
			// Everybody can list even if they have no list or get privileges. The users
			// authorization will still be used to determine what can be seen but just
			// may result in an empty list
			schema.CollectionMethods = append(schema.CollectionMethods, http.MethodGet)
		},
	})
	server.SchemaFactory.AddTemplate(schema2.Template{
		Group: "",
		Kind:  "Namespace",
		Customize: func(schema *types.APISchema) {
			// Everybody can list even if they have no list or get privileges. The users
			// authorization will still be used to determine what can be seen but just
			// may result in an empty list
			schema.CollectionMethods = append(schema.CollectionMethods, http.MethodGet)
		},
	})

	return nil
}
