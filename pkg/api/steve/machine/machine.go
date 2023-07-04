package machine

import (
	"net/http"

	"github.com/ranger/apiserver/pkg/types"
	"github.com/ranger/ranger/pkg/capr"
	"github.com/ranger/ranger/pkg/wrangler"
	schema2 "github.com/ranger/steve/pkg/schema"
	steve "github.com/ranger/steve/pkg/server"
)

func Register(server *steve.Server, clients *wrangler.Context) {
	sshClient := &sshClient{
		machines: clients.CAPI.Machine(),
		secrets:  clients.Core.Secret(),
	}

	server.SchemaFactory.AddTemplate(schema2.Template{
		Group: "cluster.x-k8s.io",
		Kind:  "Machine",
		Customize: func(schema *types.APISchema) {
			if schema.LinkHandlers == nil {
				schema.LinkHandlers = map[string]http.Handler{}
			}
			schema.LinkHandlers["shell"] = sshClient
			schema.LinkHandlers["sshkeys"] = sshClient
			schema.Formatter = func(request *types.APIRequest, resource *types.RawResource) {
				if err := request.AccessControl.CanUpdate(request, types.APIObject{}, request.Schema); err != nil ||
					resource.APIObject.Data().String("spec", "infrastructureRef", "apiVersion") != capr.RKEMachineAPIVersion {
					delete(resource.Links, "shell")
					delete(resource.Links, "sshkeys")
				}
			}
		},
	})
}
