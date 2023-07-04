package managementapi

import (
	"context"

	normanapi "github.com/ranger/norman/api"
	"github.com/ranger/ranger/pkg/auth/tokens"
	"github.com/ranger/ranger/pkg/clustermanager"
	"github.com/ranger/ranger/pkg/controllers/management/auth"
	v3cluster "github.com/ranger/ranger/pkg/controllers/management/cluster"
	podsecuritypolicy2 "github.com/ranger/ranger/pkg/controllers/management/podsecuritypolicy"
	"github.com/ranger/ranger/pkg/controllers/managementapi/catalog"
	"github.com/ranger/ranger/pkg/controllers/managementapi/dynamicschema"
	"github.com/ranger/ranger/pkg/controllers/managementapi/samlconfig"
	"github.com/ranger/ranger/pkg/controllers/managementapi/usercontrollers"
	whitelistproxyKontainerDriver "github.com/ranger/ranger/pkg/controllers/managementapi/whitelistproxy/kontainerdriver"
	whitelistproxyNodeDriver "github.com/ranger/ranger/pkg/controllers/managementapi/whitelistproxy/nodedriver"
	"github.com/ranger/ranger/pkg/controllers/managementuser/clusterauthtoken"
	"github.com/ranger/ranger/pkg/controllers/managementuser/rbac"
	"github.com/ranger/ranger/pkg/controllers/managementuser/rbac/podsecuritypolicy"
	"github.com/ranger/ranger/pkg/controllers/managementuserlegacy/monitoring"
	"github.com/ranger/ranger/pkg/types/config"
)

func Register(ctx context.Context, scaledContext *config.ScaledContext, clusterManager *clustermanager.Manager, server *normanapi.Server) error {
	if err := registerIndexers(scaledContext); err != nil {
		return err
	}

	catalog.Register(ctx, scaledContext)
	dynamicschema.Register(ctx, scaledContext, server.Schemas)
	whitelistproxyNodeDriver.Register(ctx, scaledContext)
	whitelistproxyKontainerDriver.Register(ctx, scaledContext)
	samlconfig.Register(ctx, scaledContext)
	usercontrollers.Register(ctx, scaledContext, clusterManager)
	return nil
}

func registerIndexers(scaledContext *config.ScaledContext) error {
	if err := clusterauthtoken.RegisterIndexers(scaledContext); err != nil {
		return err
	}
	if err := rbac.RegisterIndexers(scaledContext); err != nil {
		return err
	}
	if err := monitoring.RegisterIndexers(scaledContext); err != nil {
		return err
	}
	if err := auth.RegisterIndexers(scaledContext); err != nil {
		return err
	}
	if err := tokens.RegisterIndexer(scaledContext); err != nil {
		return err
	}
	if err := podsecuritypolicy.RegisterIndexers(scaledContext); err != nil {
		return err
	}
	if err := podsecuritypolicy2.RegisterIndexers(scaledContext); err != nil {
		return err
	}
	v3cluster.RegisterIndexers(scaledContext)
	return nil
}
