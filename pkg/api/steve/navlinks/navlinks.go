package navlinks

import (
	"context"

	"github.com/ranger/apiserver/pkg/types"
	schema2 "github.com/ranger/steve/pkg/schema"
	steve "github.com/ranger/steve/pkg/server"
)

func Register(ctx context.Context, server *steve.Server) {
	server.SchemaFactory.AddTemplate(schema2.Template{
		Group: "ui.cattle.io",
		Kind:  "NavLink",
		StoreFactory: func(innerStore types.Store) types.Store {
			return &store{
				Store: innerStore,
			}
		},
	})
}
