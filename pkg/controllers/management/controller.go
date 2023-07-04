package management

import (
	"context"

	"github.com/ranger/ranger/pkg/clustermanager"
	"github.com/ranger/ranger/pkg/controllers/management/agentupgrade"
	"github.com/ranger/ranger/pkg/controllers/management/auth"
	"github.com/ranger/ranger/pkg/controllers/management/certsexpiration"
	"github.com/ranger/ranger/pkg/controllers/management/cloudcredential"
	"github.com/ranger/ranger/pkg/controllers/management/cluster"
	"github.com/ranger/ranger/pkg/controllers/management/clusterdeploy"
	"github.com/ranger/ranger/pkg/controllers/management/clustergc"
	"github.com/ranger/ranger/pkg/controllers/management/clusterprovisioner"
	"github.com/ranger/ranger/pkg/controllers/management/clusterstats"
	"github.com/ranger/ranger/pkg/controllers/management/clusterstatus"
	"github.com/ranger/ranger/pkg/controllers/management/clustertemplate"
	"github.com/ranger/ranger/pkg/controllers/management/drivers/kontainerdriver"
	"github.com/ranger/ranger/pkg/controllers/management/drivers/nodedriver"
	"github.com/ranger/ranger/pkg/controllers/management/etcdbackup"
	"github.com/ranger/ranger/pkg/controllers/management/kontainerdrivermetadata"
	"github.com/ranger/ranger/pkg/controllers/management/node"
	"github.com/ranger/ranger/pkg/controllers/management/nodepool"
	"github.com/ranger/ranger/pkg/controllers/management/nodetemplate"
	"github.com/ranger/ranger/pkg/controllers/management/podsecuritypolicy"
	"github.com/ranger/ranger/pkg/controllers/management/rbac"
	"github.com/ranger/ranger/pkg/controllers/management/restrictedadminrbac"
	"github.com/ranger/ranger/pkg/controllers/management/rkeworkerupgrader"
	"github.com/ranger/ranger/pkg/controllers/management/secretmigrator"
	"github.com/ranger/ranger/pkg/controllers/management/settings"
	"github.com/ranger/ranger/pkg/controllers/management/usercontrollers"
	"github.com/ranger/ranger/pkg/controllers/managementlegacy"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/wrangler"
)

func Register(ctx context.Context, management *config.ManagementContext, manager *clustermanager.Manager, wrangler *wrangler.Context) {
	// auth handlers need to run early to create namespaces that back clusters and projects
	// also, these handlers are purely in the mgmt plane, so they are lightweight compared to those that interact with machines and clusters
	auth.RegisterEarly(ctx, management, manager)
	usercontrollers.RegisterEarly(ctx, management, manager)

	// a-z
	agentupgrade.Register(ctx, management)
	certsexpiration.Register(ctx, management)
	cluster.Register(ctx, management)
	clusterdeploy.Register(ctx, management, manager)
	clustergc.Register(ctx, management)
	clusterprovisioner.Register(ctx, management)
	clusterstats.Register(ctx, management, manager)
	clusterstatus.Register(ctx, management)
	kontainerdriver.Register(ctx, management)
	kontainerdrivermetadata.Register(ctx, management)
	nodedriver.Register(ctx, management)
	nodepool.Register(ctx, management)
	cloudcredential.Register(ctx, management)
	node.Register(ctx, management, manager)
	podsecuritypolicy.Register(ctx, management)
	etcdbackup.Register(ctx, management)
	clustertemplate.Register(ctx, management)
	nodetemplate.Register(ctx, management)
	rkeworkerupgrader.Register(ctx, management, manager.ScaledContext)
	rbac.Register(ctx, management)
	restrictedadminrbac.Register(ctx, management, wrangler)
	secretmigrator.Register(ctx, management)
	settings.Register(ctx, management)
	managementlegacy.Register(ctx, management, manager)

	// Ensure caches are available for user controllers, these are used as part of
	// registration
	management.Management.ClusterAlertGroups("").Controller()
	management.Management.ClusterAlertRules("").Controller()

	// Register last
	auth.RegisterLate(ctx, management)
}
