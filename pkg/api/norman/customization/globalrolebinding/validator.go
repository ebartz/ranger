package globalrolebinding

import (
	"net/http"

	"github.com/ranger/norman/httperror"
	"github.com/ranger/norman/types"
)

func Validator(request *types.APIContext, schema *types.Schema, data map[string]interface{}) error {
	if request.Method == http.MethodPut {
		return nil
	}

	hasSingleSubject := (data["groupPrincipalId"] != nil && data["userId"] == nil) ||
		(data["groupPrincipalId"] == nil && data["userId"] != nil)

	if !hasSingleSubject {
		return httperror.NewAPIError(httperror.InvalidBodyContent, "must contain field [groupPrincipalId] "+
			"OR field [userId]")
	}
	return nil
}
