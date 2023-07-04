package providers

import (
	"context"

	"github.com/ranger/norman/store/subtype"
	"github.com/ranger/norman/types"
	"github.com/ranger/ranger/pkg/auth/api/secrets"
	client "github.com/ranger/ranger/pkg/client/generated/management/v3"
	"github.com/ranger/ranger/pkg/namespace"
	managementschema "github.com/ranger/ranger/pkg/schemas/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/types/config"
)

var authConfigTypes = []string{
	client.GithubConfigType,
	client.LocalConfigType,
	client.ActiveDirectoryConfigType,
	client.AzureADConfigType,
	client.OpenLdapConfigType,
	client.FreeIpaConfigType,
	client.PingConfigType,
	client.ADFSConfigType,
	client.KeyCloakConfigType,
	client.OKTAConfigType,
	client.ShibbolethConfigType,
	client.GoogleOauthConfigType,
	client.OIDCConfigType,
	client.KeyCloakOIDCConfigType,
}

func SetupAuthConfig(ctx context.Context, management *config.ScaledContext, schemas *types.Schemas) {
	Configure(ctx, management)

	authConfigBaseSchema := schemas.Schema(&managementschema.Version, client.AuthConfigType)
	authConfigBaseSchema.Store = secrets.Wrap(authConfigBaseSchema.Store, management.Core.Secrets(namespace.GlobalNamespace))
	for _, authConfigSubtype := range authConfigTypes {
		subSchema := schemas.Schema(&managementschema.Version, authConfigSubtype)
		GetProviderByType(authConfigSubtype).CustomizeSchema(subSchema)
		subSchema.Store = subtype.NewSubTypeStore(authConfigSubtype, authConfigBaseSchema.Store)
	}
}
