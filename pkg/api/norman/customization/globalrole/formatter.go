package globalrole

import (
	"github.com/ranger/norman/types"
)

func Formatter(request *types.APIContext, resource *types.RawResource) {
	if resource.Values["builtin"] == true {
		delete(resource.Links, "remove")
	}
}
