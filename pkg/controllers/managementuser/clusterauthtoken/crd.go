package clusterauthtoken

import (
	"context"

	"github.com/ranger/norman/store/crd"
	"github.com/ranger/ranger/pkg/api/scheme"
	client "github.com/ranger/ranger/pkg/client/generated/cluster/v3"
	clusterSchema "github.com/ranger/ranger/pkg/schemas/cluster.cattle.io/v3"
	"github.com/ranger/ranger/pkg/types/config"
)

func CRDSetup(ctx context.Context, apiContext *config.UserOnlyContext) error {
	factory, err := crd.NewFactoryFromClient(apiContext.RESTConfig)
	if err != nil {
		return err
	}
	factory.BatchCreateCRDs(ctx, config.UserStorageContext, scheme.Scheme, apiContext.Schemas, &clusterSchema.Version,
		client.ClusterAuthTokenType,
		client.ClusterUserAttributeType,
	)
	return factory.BatchWait()
}
