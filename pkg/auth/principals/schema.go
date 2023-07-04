package principals

import (
	"context"
	"net/url"

	"github.com/ranger/norman/types"
	"github.com/ranger/ranger/pkg/auth/requests"
	client "github.com/ranger/ranger/pkg/client/generated/management/v3"
	managementSchema "github.com/ranger/ranger/pkg/schemas/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/types/config"
)

func Schema(ctx context.Context, clusterRouter requests.ClusterRouter, management *config.ScaledContext, schemas *types.Schemas) error {
	p := newPrincipalsHandler(ctx, clusterRouter, management)
	schema := schemas.Schema(&managementSchema.Version, client.PrincipalType)
	schema.ActionHandler = p.actions
	schema.ListHandler = p.list
	schema.CollectionFormatter = collectionFormatter
	schema.Formatter = formatter
	return nil
}

func collectionFormatter(apiContext *types.APIContext, collection *types.GenericCollection) {
	collection.AddAction(apiContext, "search")
}

func formatter(request *types.APIContext, resource *types.RawResource) {
	schema := request.Schemas.Schema(&managementSchema.Version, client.PrincipalType)
	resource.Links = map[string]string{"self": request.URLBuilder.ResourceLinkByID(schema, url.PathEscape(resource.ID))}
}
