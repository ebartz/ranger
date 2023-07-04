package steve

import (
	"context"

	"github.com/ranger/ranger/pkg/api/steve/catalog"
	"github.com/ranger/ranger/pkg/api/steve/clusters"
	"github.com/ranger/ranger/pkg/api/steve/disallow"
	"github.com/ranger/ranger/pkg/api/steve/machine"
	"github.com/ranger/ranger/pkg/api/steve/navlinks"
	"github.com/ranger/ranger/pkg/api/steve/settings"
	"github.com/ranger/ranger/pkg/api/steve/userpreferences"
	"github.com/ranger/ranger/pkg/wrangler"
	steve "github.com/ranger/steve/pkg/server"
)

func Setup(ctx context.Context, server *steve.Server, config *wrangler.Context) error {
	userpreferences.Register(server.BaseSchemas, server.ClientFactory)
	if err := clusters.Register(ctx, server, config); err != nil {
		return err
	}
	machine.Register(server, config)
	navlinks.Register(ctx, server)
	settings.Register(server)
	disallow.Register(server)
	return catalog.Register(ctx,
		server,
		config.HelmOperations,
		config.CatalogContentManager)
}
