package factory

import (
	"github.com/ranger/norman/types"
	"github.com/ranger/norman/types/factory"
	m "github.com/ranger/norman/types/mapper"
	"github.com/ranger/ranger/pkg/schemas/mapper"
)

func Schemas(version *types.APIVersion) *types.Schemas {
	schemas := factory.Schemas(version)
	baseFunc := schemas.DefaultMappers
	schemas.DefaultMappers = func() []types.Mapper {
		mappers := append([]types.Mapper{
			&mapper.Status{},
		}, baseFunc()...)
		mappers = append(mappers, &m.Scope{
			If: types.NamespaceScope,
			Mappers: []types.Mapper{
				&mapper.NamespaceIDMapper{},
			},
		}, &mapper.NamespaceReference{
			VersionPath: "/v3/project",
		})
		return mappers
	}
	basePostFunc := schemas.DefaultPostMappers
	schemas.DefaultPostMappers = func() []types.Mapper {
		return append(basePostFunc(), &mapper.Creator{})
	}
	return schemas
}
