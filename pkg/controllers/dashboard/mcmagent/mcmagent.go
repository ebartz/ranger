package mcmagent

import (
	"context"

	"github.com/ranger/ranger/pkg/controllers/managementagent"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/wrangler"
)

func Register(ctx context.Context, wrangler *wrangler.Context) error {
	userContext, err := config.NewUserOnlyContext(wrangler)
	if err != nil {
		return err
	}
	return managementagent.Register(ctx, userContext)
}
