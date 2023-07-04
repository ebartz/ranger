package etcdbackup

import (
	"github.com/ranger/norman/types"
	"github.com/ranger/norman/types/convert"
	"github.com/ranger/norman/types/values"
)

func Formatter(apiContext *types.APIContext, resource *types.RawResource) {
	state := convert.ToString(resource.Values["state"])
	if state == "activating" {
		for _, cond := range convert.ToMapSlice(values.GetValueN(resource.Values, "status", "conditions")) {
			if cond["type"] == "Completed" {
				if cond["status"] == "False" && convert.ToString(cond["reason"]) == "Error" {
					resource.Values["state"] = "failed"
				}
				break
			}
		}
	}
}
