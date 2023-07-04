package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ranger/norman/types"
	"github.com/ranger/ranger/pkg/api/scheme"
	v3 "github.com/ranger/ranger/pkg/apis/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/ref"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	management "github.com/ranger/ranger/tests/framework/clients/ranger/generated/management/v3"
	"github.com/ranger/ranger/tests/framework/extensions/defaults"
	"github.com/ranger/ranger/tests/framework/extensions/kubeapi/rbac"
	kubeapiSecrets "github.com/ranger/ranger/tests/framework/extensions/kubeapi/secrets"
	"github.com/ranger/ranger/tests/framework/extensions/secrets"
	"github.com/ranger/ranger/tests/framework/pkg/wait"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	kwait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	rtbOwnerLabel = "authz.cluster.cattle.io/rtb-owner-updated"
)

var timeout = int64(60 * 3)

// CreateUserWithRole is helper function that creates a user with a role or multiple roles
func CreateUserWithRole(rangerClient *ranger.Client, user *management.User, roles ...string) (*management.User, error) {
	createdUser, err := rangerClient.Management.User.Create(user)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		roleBinding := &management.GlobalRoleBinding{
			GlobalRoleID: role,
			UserID:       createdUser.ID,
		}

		_, err = rangerClient.Management.GlobalRoleBinding.Create(roleBinding)
		if err != nil {
			return nil, err
		}
	}

	return createdUser, nil
}

// AddProjectMember is a helper function that adds a project role to `user`. It uses the watch.WatchWait ensure BackingNamespaceCreated is true
func AddProjectMember(rangerClient *ranger.Client, project *management.Project, user *management.User, projectRole string) error {
	role := &management.ProjectRoleTemplateBinding{
		ProjectID:       project.ID,
		UserPrincipalID: user.PrincipalIDs[0],
		RoleTemplateID:  projectRole,
	}

	name := strings.Split(project.ID, ":")[1]

	adminClient, err := ranger.NewClient(rangerClient.RangerConfig.AdminToken, rangerClient.Session)
	if err != nil {
		return err
	}

	opts := metav1.ListOptions{
		FieldSelector:  "metadata.name=" + name,
		TimeoutSeconds: &timeout,
	}
	watchInterface, err := adminClient.GetManagementWatchInterface(management.ProjectType, opts)
	if err != nil {
		return err
	}

	checkFunc := func(event watch.Event) (ready bool, err error) {
		projectUnstructured := event.Object.(*unstructured.Unstructured)
		project := &v3.Project{}
		err = scheme.Scheme.Convert(projectUnstructured, project, projectUnstructured.GroupVersionKind())
		if err != nil {
			return false, err
		}
		if v3.NamespaceBackedResource.IsTrue(project) {
			return true, nil
		}

		return false, nil
	}

	err = wait.WatchWait(watchInterface, checkFunc)
	if err != nil {
		return err
	}

	roleTemplateResp, err := rangerClient.Management.ProjectRoleTemplateBinding.Create(role)
	if err != nil {
		return err
	}

	err = kwait.Poll(500*time.Millisecond, 2*time.Minute, func() (done bool, err error) {
		projectRoleTemplate, err := rangerClient.Management.ProjectRoleTemplateBinding.ByID(roleTemplateResp.ID)
		if err != nil {
			return false, err
		}
		if projectRoleTemplate != nil && projectRoleTemplate.UserID == user.ID && projectRoleTemplate.ProjectID == project.ID {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return err
	}

	adminDynamicClient, err := adminClient.GetDownStreamClusterClient(project.ClusterID)
	if err != nil {
		return err
	}

	steveClient, err := adminClient.Steve.ProxyDownstream(project.ClusterID)
	if err != nil {
		return err
	}

	secretOpts := metav1.ListOptions{
		FieldSelector:  "metadata.namespace=" + "cattle-impersonation-system",
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	}

	var numOfActiveSecrets int
	err = kwait.Poll(500*time.Millisecond, 2*time.Minute, func() (done bool, err error) {
		secretsList, err := adminDynamicClient.Resource(kubeapiSecrets.SecretGroupVersionResource).List(context.TODO(), secretOpts)
		if err != nil {
			return false, err
		}

		for _, secret := range secretsList.Items {

			if strings.Contains(secret.GetName(), user.ID) {
				secretID := fmt.Sprintf("%s/%s", secret.GetNamespace(), secret.GetName())
				steveSecret, err := steveClient.SteveType(secrets.SecretSteveType).ByID(secretID)
				if err != nil {
					return false, err
				}

				if steveSecret.ObjectMeta.State.Name == "active" {
					numOfActiveSecrets += 1
				}

				if numOfActiveSecrets == 2 {
					return true, nil
				}
			}
		}

		return false, nil
	})

	return err
}

// RemoveProjectMember is a helper function that removes the project role from `user`
func RemoveProjectMember(rangerClient *ranger.Client, user *management.User) error {
	roles, err := rangerClient.Management.ProjectRoleTemplateBinding.List(&types.ListOpts{})
	if err != nil {
		return err
	}

	var roleToDelete management.ProjectRoleTemplateBinding

	for _, role := range roles.Data {
		if role.UserID == user.ID {
			roleToDelete = role
			break
		}
	}

	var backoff = kwait.Backoff{
		Duration: 100 * time.Millisecond,
		Factor:   1,
		Jitter:   0,
		Steps:    5,
	}
	err = rangerClient.Management.ProjectRoleTemplateBinding.Delete(&roleToDelete)
	if err != nil {
		return err
	}
	err = kwait.ExponentialBackoff(backoff, func() (done bool, err error) {
		clusterID, projName := ref.Parse(roleToDelete.ProjectID)
		req, err := labels.NewRequirement(rtbOwnerLabel, selection.Equals, []string{fmt.Sprintf("%s_%s", projName, roleToDelete.Name)})
		if err != nil {
			return false, err
		}

		downstreamRBs, err := rbac.ListRoleBindings(rangerClient, clusterID, "", metav1.ListOptions{
			LabelSelector: labels.NewSelector().Add(*req).String(),
		})
		if err != nil {
			return false, err
		}
		if len(downstreamRBs.Items) != 0 {
			return false, nil
		}
		return true, nil
	})
	return err
}

// AddClusterRoleToUser is a helper function that adds a cluster role to `user`.
func AddClusterRoleToUser(rangerClient *ranger.Client, cluster *management.Cluster, user *management.User, clusterRole string) error {
	role := &management.ClusterRoleTemplateBinding{
		ClusterID:       cluster.Resource.ID,
		UserPrincipalID: user.PrincipalIDs[0],
		RoleTemplateID:  clusterRole,
	}

	adminClient, err := ranger.NewClient(rangerClient.RangerConfig.AdminToken, rangerClient.Session)
	if err != nil {
		return err
	}

	opts := metav1.ListOptions{
		FieldSelector:  "metadata.name=" + cluster.ID,
		TimeoutSeconds: &timeout,
	}
	watchInterface, err := adminClient.GetManagementWatchInterface(management.ClusterType, opts)
	if err != nil {
		return err
	}

	checkFunc := func(event watch.Event) (ready bool, err error) {
		clusterUnstructured := event.Object.(*unstructured.Unstructured)
		cluster := &v3.Cluster{}

		err = scheme.Scheme.Convert(clusterUnstructured, cluster, clusterUnstructured.GroupVersionKind())
		if err != nil {
			return false, err
		}
		if cluster.Annotations == nil || cluster.Annotations["field.cattle.io/creatorId"] == "" {
			// no cluster creator, no roles to populate. This will be the case for the "local" cluster.
			return true, nil
		}

		v3.ClusterConditionInitialRolesPopulated.CreateUnknownIfNotExists(cluster)
		if v3.ClusterConditionInitialRolesPopulated.IsUnknown(cluster) || v3.ClusterConditionInitialRolesPopulated.IsTrue(cluster) {
			return true, nil
		}
		return false, nil
	}

	err = wait.WatchWait(watchInterface, checkFunc)
	if err != nil {
		return err
	}

	roleTemplateResp, err := rangerClient.Management.ClusterRoleTemplateBinding.Create(role)
	if err != nil {
		return err
	}

	err = kwait.Poll(600*time.Millisecond, 3*time.Minute, func() (done bool, err error) {
		clusterRoleTemplate, err := rangerClient.Management.ClusterRoleTemplateBinding.ByID(roleTemplateResp.ID)
		if err != nil {
			return false, err
		}
		if clusterRoleTemplate != nil {
			return true, nil
		}

		return false, nil
	})

	return err

}

// RemoveClusterRoleFromUser is a helper function that removes the user from cluster
func RemoveClusterRoleFromUser(rangerClient *ranger.Client, user *management.User) error {
	roles, err := rangerClient.Management.ClusterRoleTemplateBinding.List(&types.ListOpts{})
	if err != nil {
		return err
	}

	var roleToDelete management.ClusterRoleTemplateBinding

	for _, role := range roles.Data {
		if role.UserID == user.ID {
			roleToDelete = role
			break
		}
	}

	if err = rangerClient.Management.ClusterRoleTemplateBinding.Delete(&roleToDelete); err != nil {
		return err
	}

	var backoff = kwait.Backoff{
		Duration: 100 * time.Millisecond,
		Factor:   1,
		Jitter:   0,
		Steps:    5,
	}

	err = kwait.ExponentialBackoff(backoff, func() (done bool, err error) {
		req, err := labels.NewRequirement(rtbOwnerLabel, selection.Equals, []string{fmt.Sprintf("%s_%s", roleToDelete.ClusterID, roleToDelete.Name)})
		if err != nil {
			return false, err
		}

		downstreamCRBs, err := rbac.ListClusterRoleBindings(rangerClient, roleToDelete.ClusterID, metav1.ListOptions{
			LabelSelector: labels.NewSelector().Add(*req).String(),
		})
		if err != nil {
			return false, err
		}
		if len(downstreamCRBs.Items) != 0 {
			return false, nil
		}
		return true, nil
	})
	return err
}

// GetUserIDByName is a helper function that returns the user ID by name
func GetUserIDByName(client *ranger.Client, username string) (string, error) {
	userList, err := client.Management.User.List(&types.ListOpts{})
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	for _, user := range userList.Data {
		if user.Username == username {
			return user.ID, nil
		}
	}

	return "", nil
}
