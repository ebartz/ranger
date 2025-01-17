package namespace

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ranger/norman/api/access"
	"github.com/ranger/norman/httperror"
	"github.com/ranger/norman/parse"
	"github.com/ranger/norman/types"
	"github.com/ranger/norman/types/convert"
	client "github.com/ranger/ranger/pkg/client/generated/cluster/v3"
	"github.com/ranger/ranger/pkg/clustermanager"
	"github.com/ranger/ranger/pkg/controllers/managementagent/nslabels"
	"github.com/ranger/ranger/pkg/controllers/managementuserlegacy/helm"
	"github.com/ranger/ranger/pkg/rbac"
	"github.com/ranger/ranger/pkg/ref"
	schema "github.com/ranger/ranger/pkg/schemas/cluster.cattle.io/v3"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/cache"
)

var (
	namespaceOwnerMap = cache.NewLRUExpireCache(1000)
)

func updateNamespaceOwnerMap(apiContext *types.APIContext) error {
	var namespaces []client.Namespace
	if err := access.List(apiContext, &schema.Version, client.NamespaceType, &types.QueryOptions{}, &namespaces); err != nil {
		return err
	}

	for _, namespace := range namespaces {
		namespaceOwnerMap.Add(namespace.Name, namespace.ProjectID, time.Hour)
	}

	return nil
}

func ProjectMap(apiContext *types.APIContext, refresh bool) (map[string]string, error) {
	if refresh {
		err := updateNamespaceOwnerMap(apiContext)
		if err != nil {
			return nil, err
		}
	}

	data := map[string]string{}
	for _, key := range namespaceOwnerMap.Keys() {
		if val, ok := namespaceOwnerMap.Get(key); ok {
			data[key.(string)] = val.(string)
		}
	}

	return data, nil
}

type ActionWrapper struct {
	ClusterManager *clustermanager.Manager
}

func (w ActionWrapper) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	actionInput, err := parse.ReadBody(apiContext.Request)
	if err != nil {
		return err
	}

	if !canUpdateNS(apiContext, nil) {
		return httperror.NewAPIError(httperror.NotFound, "not found")
	}

	switch actionName {
	case "move":
		clusterID := w.ClusterManager.ClusterName(apiContext)
		_, projectID := ref.Parse(convert.ToString(actionInput["projectId"]))
		userContext, err := w.ClusterManager.UserContextNoControllers(clusterID)
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return err
			}
			return httperror.NewAPIError(httperror.NotFound, err.Error())
		}
		if projectID != "" {
			project, err := userContext.Management.Management.Projects(clusterID).Get(projectID, metav1.GetOptions{})
			if err != nil {
				return err
			}
			if project.Spec.ResourceQuota != nil {
				return errors.Errorf("can't move namespace. Project %s has resource quota set", project.Spec.DisplayName)
			}
		}
		nsClient := userContext.Core.Namespaces("")
		ns, err := nsClient.Get(apiContext.ID, metav1.GetOptions{})
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return err
			}
			return httperror.NewAPIError(httperror.NotFound, err.Error())
		}
		if ns.Annotations[helm.AppIDsLabel] != "" {
			return errors.New("namespace is currently being used")
		}
		if projectID == "" {
			delete(ns.Annotations, nslabels.ProjectIDFieldLabel)
			delete(ns.Labels, nslabels.ProjectIDFieldLabel)
		} else {
			ns.Annotations[nslabels.ProjectIDFieldLabel] = convert.ToString(actionInput["projectId"])
		}
		if _, err := nsClient.Update(ns); err != nil {
			return err
		}
	default:
		return errors.New("invalid action")
	}
	return nil
}

func NewFormatter(next types.Formatter) types.Formatter {
	return func(request *types.APIContext, resource *types.RawResource) {
		if next != nil {
			next(request, resource)
		}
		annotations := convert.ToMapInterface(resource.Values["annotations"])
		canUpdate := canUpdateNS(request, resource)
		if canUpdate && convert.ToString(annotations[helm.AppIDsLabel]) == "" {
			resource.AddAction(request, "move")
		}
	}
}

func canUpdateNS(apiContext *types.APIContext, resource *types.RawResource) bool {
	obj := rbac.ObjFromContext(apiContext, resource)
	// the user must have * permissions on namespace, the create-ns role alone won't return true here
	return apiContext.AccessControl.CanDo("", "namespaces", "update", apiContext, obj, apiContext.Schema) == nil
}
