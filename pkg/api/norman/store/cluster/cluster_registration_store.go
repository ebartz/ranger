package cluster

import (
	"github.com/ranger/norman/types"
	"github.com/ranger/wrangler/pkg/randomtoken"
)

type RegistrationTokenStore struct {
	types.Store
}

func (r *RegistrationTokenStore) Create(apiContext *types.APIContext, schema *types.Schema, data map[string]interface{}) (map[string]interface{}, error) {
	if data != nil {
		token, err := randomtoken.Generate()
		if err != nil {
			return nil, err
		}
		data["token"] = token
	}

	return r.Store.Create(apiContext, schema, data)
}
