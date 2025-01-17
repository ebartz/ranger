package userstored

import (
	"context"
	"strings"

	"github.com/ranger/norman/store/proxy"
	"github.com/ranger/norman/types"
	"github.com/ranger/ranger/pkg/api/scheme"
	clusterSchema "github.com/ranger/ranger/pkg/schemas/cluster.cattle.io/v3"
	schema "github.com/ranger/ranger/pkg/schemas/project.cattle.io/v3"
	"github.com/ranger/ranger/pkg/types/config"
)

type storeWrapperFunc func(types.Store) types.Store

func addProxyStore(ctx context.Context, schemas *types.Schemas, context *config.ScaledContext, schemaType, apiVersion string, storeWrapper storeWrapperFunc) *types.Schema {
	s := schemas.Schema(&schema.Version, schemaType)
	if s == nil {
		s = schemas.Schema(&clusterSchema.Version, schemaType)
	}

	if s == nil {
		panic("Failed to find schema " + schemaType)
	}

	prefix := []string{"api"}
	kind := s.CodeName
	plural := strings.ToLower(s.PluralName)

	var version, group string
	parts := strings.SplitN(apiVersion, "/", 2)
	if len(parts) == 1 {
		version = parts[0]
	} else {
		group = parts[0]
		version = parts[1]
		prefix = []string{"apis"}
	}

	s.Store = proxy.NewProxyStore(ctx, context.ClientGetter,
		config.UserStorageContext,
		scheme.Scheme,
		prefix,
		group,
		version,
		kind,
		plural)

	if storeWrapper != nil {
		s.Store = storeWrapper(s.Store)
	}

	return s
}
