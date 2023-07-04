package mapper

import (
	"github.com/ranger/norman/types"
	"github.com/ranger/norman/types/convert"
	"github.com/ranger/norman/types/values"
)

type ContainerProbeHandler struct {
}

func (n ContainerProbeHandler) FromInternal(data map[string]interface{}) {
	value := values.GetValueN(data, "tcpSocket", "port")
	if !convert.IsAPIObjectEmpty(value) {
		data["tcp"] = true
	}
}

func (n ContainerProbeHandler) ToInternal(data map[string]interface{}) error {
	return nil
}

func (n ContainerProbeHandler) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
