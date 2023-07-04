package preference

import (
	"strings"

	"github.com/ranger/norman/store/transform"
	"github.com/ranger/norman/types"
	"github.com/ranger/norman/types/convert"
	"github.com/ranger/ranger/pkg/api/norman/store/userscope"
	client "github.com/ranger/ranger/pkg/client/generated/management/v3"
	v1 "github.com/ranger/ranger/pkg/generated/norman/core/v1"
)

const (
	NamespaceID = client.PreferenceFieldNamespaceId
)

func NewStore(nsClient v1.NamespaceInterface, store types.Store) types.Store {
	return userscope.NewStore(nsClient,
		&transform.Store{
			Store:       store,
			Transformer: transformer,
		})
}

func transformer(apiContext *types.APIContext, schema *types.Schema, data map[string]interface{}, opts *types.QueryOptions) (map[string]interface{}, error) {
	if data == nil {
		return nil, nil
	}

	ns := convert.ToString(data[NamespaceID])
	id := convert.ToString(data[types.ResourceFieldID])

	id = strings.TrimPrefix(id, ns+":")

	data[client.PreferenceFieldName] = id
	data[types.ResourceFieldID] = id

	return data, nil
}
