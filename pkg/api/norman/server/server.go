package server

import (
	"context"
	"net/http"

	responsewriter "github.com/ranger/apiserver/pkg/middleware"
	"github.com/ranger/norman/api/builtin"
	"github.com/ranger/norman/pkg/subscribe"
	"github.com/ranger/ranger/pkg/api/norman"
	"github.com/ranger/ranger/pkg/api/norman/server/managementstored"
	"github.com/ranger/ranger/pkg/api/norman/server/userstored"
	"github.com/ranger/ranger/pkg/clustermanager"
	"github.com/ranger/ranger/pkg/controllers/managementapi"
	clusterSchema "github.com/ranger/ranger/pkg/schemas/cluster.cattle.io/v3"
	managementSchema "github.com/ranger/ranger/pkg/schemas/management.cattle.io/v3"
	projectSchema "github.com/ranger/ranger/pkg/schemas/project.cattle.io/v3"
	"github.com/ranger/ranger/pkg/types/config"
)

func New(ctx context.Context, scaledContext *config.ScaledContext, clusterManager *clustermanager.Manager,
	k8sProxy http.Handler, localClusterEnabled bool) (http.Handler, error) {
	subscribe.Register(&builtin.Version, scaledContext.Schemas)
	subscribe.Register(&managementSchema.Version, scaledContext.Schemas)
	subscribe.Register(&clusterSchema.Version, scaledContext.Schemas)
	subscribe.Register(&projectSchema.Version, scaledContext.Schemas)

	if err := managementstored.Setup(ctx, scaledContext, clusterManager, k8sProxy, localClusterEnabled); err != nil {
		return nil, err
	}

	if err := userstored.Setup(ctx, scaledContext, clusterManager, k8sProxy); err != nil {
		return nil, err
	}

	server, err := norman.NewServer(scaledContext.Schemas)
	if err != nil {
		return nil, err
	}
	server.AccessControl = scaledContext.AccessControl

	if err := managementapi.Register(ctx, scaledContext, clusterManager, server); err != nil {
		return nil, err
	}

	chainGzip := responsewriter.Chain{responsewriter.Gzip, responsewriter.ContentType}
	return chainGzip.Handler(server), nil
}
