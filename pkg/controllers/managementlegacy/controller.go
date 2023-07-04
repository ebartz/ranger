package managementlegacy

import (
	"context"

	"github.com/ranger/ranger/pkg/clustermanager"
	"github.com/ranger/ranger/pkg/controllers/managementlegacy/catalog"
	"github.com/ranger/ranger/pkg/controllers/managementlegacy/compose"
	"github.com/ranger/ranger/pkg/controllers/managementlegacy/globaldns"
	"github.com/ranger/ranger/pkg/controllers/managementlegacy/multiclusterapp"
	"github.com/ranger/ranger/pkg/types/config"
)

func Register(ctx context.Context, management *config.ManagementContext, manager *clustermanager.Manager) {
	catalog.Register(ctx, management)
	compose.Register(ctx, management, manager)
	globaldns.Register(ctx, management)
	multiclusterapp.Register(ctx, management, manager)
}
