package globaldns

import (
	"context"
	"fmt"
	"strings"

	v32 "github.com/ranger/ranger/pkg/apis/project.cattle.io/v3"

	"github.com/ranger/norman/types/convert"
	passwordutil "github.com/ranger/ranger/pkg/api/norman/store/password"
	"github.com/ranger/ranger/pkg/catalog/manager"
	cutils "github.com/ranger/ranger/pkg/catalog/utils"
	"github.com/ranger/ranger/pkg/controllers/management/rbac"
	v1 "github.com/ranger/ranger/pkg/generated/norman/core/v1"
	v3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
	pv3 "github.com/ranger/ranger/pkg/generated/norman/project.cattle.io/v3"
	"github.com/ranger/ranger/pkg/namespace"
	"github.com/ranger/ranger/pkg/project"
	"github.com/ranger/ranger/pkg/settings"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/user"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

const (
	GlobaldnsProviderCatalogLauncher = "mgmt-global-dns-provider-catalog-launcher"
	cattleCreatorIDAnnotationKey     = "field.cattle.io/creatorId"
	localClusterName                 = "local"
	templateName                     = "ranger-external-dns"
)

type ProviderCatalogLauncher struct {
	managementContext *config.ManagementContext
	Apps              pv3.AppInterface
	ProjectLister     v3.ProjectLister
	appLister         pv3.AppLister
	userManager       user.Manager
	secrets           v1.SecretInterface
	templateLister    v3.CatalogTemplateLister
	catalogManager    manager.CatalogManager
}

func newGlobalDNSProviderCatalogLauncher(ctx context.Context, mgmt *config.ManagementContext) *ProviderCatalogLauncher {
	n := &ProviderCatalogLauncher{
		managementContext: mgmt,
		Apps:              mgmt.Project.Apps(""),
		ProjectLister:     mgmt.Management.Projects("").Controller().Lister(),
		appLister:         mgmt.Project.Apps("").Controller().Lister(),
		userManager:       mgmt.UserManager,
		secrets:           mgmt.Core.Secrets(""),
		templateLister:    mgmt.Management.CatalogTemplates(metav1.NamespaceAll).Controller().Lister(),
		catalogManager:    mgmt.CatalogManager,
	}
	return n
}

// sync is called periodically and on real updates
func (n *ProviderCatalogLauncher) sync(key string, obj *v3.GlobalDnsProvider) (runtime.Object, error) {
	if obj == nil || obj.DeletionTimestamp != nil {
		return nil, n.deleteExternalDNSApp(key)
	}
	metaAccessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	creatorID, ok := metaAccessor.GetAnnotations()[rbac.CreatorIDAnn]
	if !ok {
		return nil, fmt.Errorf("GlobalDNS %v has no creatorId annotation", metaAccessor.GetName())
	}

	if err := rbac.CreateRoleAndRoleBinding(rbac.GlobalDNSProviderResource, v3.GlobalDnsProviderGroupVersionKind.Kind, obj.Name, namespace.GlobalNamespace,
		rbac.RangerManagementAPIVersion, creatorID, []string{rbac.RangerManagementAPIGroup},
		obj.UID, obj.Spec.Members, n.managementContext); err != nil {
		return nil, err
	}

	//handle external-dns deployment
	if obj.Spec.Route53ProviderConfig != nil {
		return n.handleRoute53Provider(obj)
	}

	if obj.Spec.CloudflareProviderConfig != nil {
		return n.handleCloudflareProvider(obj)
	}

	if obj.Spec.AlidnsProviderConfig != nil {
		return n.handleAlidnsProvider(obj)
	}

	return nil, nil
}

func (n *ProviderCatalogLauncher) handleRoute53Provider(obj *v3.GlobalDnsProvider) (runtime.Object, error) {
	rangerInstallUUID := settings.InstallUUID.Get()
	//create external-dns route53 provider

	secretKey := obj.Spec.Route53ProviderConfig.SecretKey
	//read the secret if found
	var err error
	if strings.HasPrefix(obj.Spec.Route53ProviderConfig.SecretKey, namespace.GlobalNamespace) {
		secretKey, err = passwordutil.GetValueForPasswordField(obj.Spec.Route53ProviderConfig.SecretKey, n.secrets)
		if err != nil {
			return nil, err
		}
	}

	answers := map[string]string{
		"provider":            "aws",
		"aws.zoneType":        obj.Spec.Route53ProviderConfig.ZoneType,
		"aws.accessKey":       obj.Spec.Route53ProviderConfig.AccessKey,
		"aws.secretKey":       secretKey,
		"txtOwnerId":          rangerInstallUUID + "_" + obj.Name,
		"rbac.create":         "true",
		"policy":              "sync",
		"aws.credentialsPath": obj.Spec.Route53ProviderConfig.CredentialsPath,
		"aws.roleArn":         obj.Spec.Route53ProviderConfig.RoleArn,
		"aws.region":          obj.Spec.Route53ProviderConfig.Region,
	}
	for k, v := range obj.Spec.Route53ProviderConfig.AdditionalOptions {
		if _, ok := answers[k]; !ok {
			answers[k] = v
		}
	}

	if obj.Spec.RootDomain != "" {
		answers["domainFilters[0]"] = obj.Spec.RootDomain
	}
	return n.createUpdateExternalDNSApp(obj, answers)
}

func (n *ProviderCatalogLauncher) handleCloudflareProvider(obj *v3.GlobalDnsProvider) (runtime.Object, error) {
	rangerInstallUUID := settings.InstallUUID.Get()

	isProxy := "true"
	if obj.Spec.CloudflareProviderConfig.ProxySetting != nil {
		isProxy = convert.ToString(*obj.Spec.CloudflareProviderConfig.ProxySetting)
	}

	secretAPIKey := obj.Spec.CloudflareProviderConfig.APIKey
	//read the secret if found
	var err error
	if strings.HasPrefix(obj.Spec.CloudflareProviderConfig.APIKey, namespace.GlobalNamespace) {
		secretAPIKey, err = passwordutil.GetValueForPasswordField(obj.Spec.CloudflareProviderConfig.APIKey, n.secrets)
		if err != nil {
			return nil, err
		}
	}

	//create external-dns route53 provider
	answers := map[string]string{
		"provider":           "cloudflare",
		"cloudflare.apiKey":  secretAPIKey,
		"cloudflare.email":   obj.Spec.CloudflareProviderConfig.APIEmail,
		"txtOwnerId":         rangerInstallUUID + "_" + obj.Name,
		"rbac.create":        "true",
		"policy":             "sync",
		"cloudflare.proxied": isProxy,
	}
	for k, v := range obj.Spec.CloudflareProviderConfig.AdditionalOptions {
		if _, ok := answers[k]; !ok {
			answers[k] = v
		}
	}

	if obj.Spec.RootDomain != "" {
		answers["domainFilters[0]"] = obj.Spec.RootDomain
	}

	return n.createUpdateExternalDNSApp(obj, answers)
}

func (n *ProviderCatalogLauncher) handleAlidnsProvider(obj *v3.GlobalDnsProvider) (runtime.Object, error) {
	rangerInstallUUID := settings.InstallUUID.Get()

	secretKey := obj.Spec.AlidnsProviderConfig.SecretKey
	//read the secret if found
	var err error
	if strings.HasPrefix(obj.Spec.AlidnsProviderConfig.SecretKey, namespace.GlobalNamespace) {
		secretKey, err = passwordutil.GetValueForPasswordField(obj.Spec.AlidnsProviderConfig.SecretKey, n.secrets)
		if err != nil {
			return nil, err
		}
	}

	//create external-dns alidns provider
	answers := map[string]string{
		"provider":               "alibabacloud",
		"alibabacloud.zoneType":  "public",
		"alibabacloud.accessKey": obj.Spec.AlidnsProviderConfig.AccessKey,
		"alibabacloud.secretKey": secretKey,
		"txtOwnerId":             rangerInstallUUID + "_" + obj.Name,
		"rbac.create":            "true",
		"policy":                 "sync",
	}
	for k, v := range obj.Spec.AlidnsProviderConfig.AdditionalOptions {
		if _, ok := answers[k]; !ok {
			answers[k] = v
		}
	}

	if obj.Spec.RootDomain != "" {
		answers["domainFilters[0]"] = obj.Spec.RootDomain
	}

	return n.createUpdateExternalDNSApp(obj, answers)
}

func (n *ProviderCatalogLauncher) createUpdateExternalDNSApp(obj *v3.GlobalDnsProvider, answers map[string]string) (runtime.Object, error) {
	//check if provider already running for this GlobalDNSProvider.
	existingApp, err := n.getProviderIfAlreadyRunning(obj.Name)
	if err != nil {
		return nil, err
	}

	if existingApp != nil {
		var shouldUpdate bool
		var appToUpdate *pv3.App

		//check if answers should be updated
		if answersDiffer(existingApp.Spec.Answers, answers) {
			appToUpdate = existingApp.DeepCopy()
			updateAnswers(appToUpdate.Spec.Answers, answers)
			shouldUpdate = true
		}

		//check if ownerRef to GlobalDnsProvider is present; drop if yes
		if len(existingApp.OwnerReferences) > 0 {
			for _, ref := range existingApp.OwnerReferences {
				if ref.UID == obj.UID {
					if appToUpdate == nil {
						appToUpdate = existingApp.DeepCopy()
					}
					removeOwnerRef(appToUpdate, ref)
					shouldUpdate = true
					break
				}
			}
		}

		if shouldUpdate {
			_, err = n.Apps.Update(appToUpdate)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				return nil, err
			}
		}
	} else {
		//create new app
		appCatalogID, err := n.getExternalDNSCatalogID(localClusterName)
		if err != nil {
			return nil, err
		}
		sysProject, err := n.getSystemProjectID()
		if err != nil {
			return nil, err
		}

		creator, err := n.userManager.EnsureUser(fmt.Sprintf("system://%s", localClusterName), "System account for Cluster "+localClusterName)
		if err != nil {
			return nil, err
		}
		toCreate := pv3.App{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{cattleCreatorIDAnnotationKey: creator.Name},
				Name:        fmt.Sprintf("%s-%s", "systemapp", obj.Name),
				Namespace:   sysProject,
			},
			Spec: v32.AppSpec{
				ProjectName:     localClusterName + ":" + sysProject,
				TargetNamespace: namespace.GlobalNamespace,
				ExternalID:      appCatalogID,
				Answers:         answers,
			},
		}
		// Now create the App instance
		_, err = n.Apps.Create(&toCreate)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return nil, err
		}
	}
	return nil, nil
}

func (n *ProviderCatalogLauncher) deleteExternalDNSApp(key string) error {
	_, globalDNSproviderName, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	//check if provider already running for this GlobalDNSProvider.
	existingApp, err := n.getProviderIfAlreadyRunning(globalDNSproviderName)
	if err != nil {
		return err
	}
	if existingApp != nil {
		//delete this app since GlobalDNSProvider is deleted
		if err := n.Apps.DeleteNamespaced(existingApp.Namespace, existingApp.Name, &metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (n *ProviderCatalogLauncher) getProviderIfAlreadyRunning(globalDNSProviderName string) (*pv3.App, error) {
	sysProject, err := n.getSystemProjectID()
	if err != nil {
		return nil, err
	}
	existingApp, err := n.appLister.Get(sysProject, fmt.Sprintf("%s-%s", "systemapp", globalDNSProviderName))

	if (err != nil && k8serrors.IsNotFound(err)) || existingApp == nil {
		return nil, nil
	} else if err != nil && !k8serrors.IsNotFound(err) {
		logrus.Errorf("GlobaldnsProviderCatalogLauncher: Error listing external-dns %v app %v", globalDNSProviderName, err)
		return nil, err
	}

	return existingApp, nil
}

func (n *ProviderCatalogLauncher) getSystemProjectID() (string, error) {
	systemProject, err := project.GetSystemProject(localClusterName, n.ProjectLister)
	if err != nil {
		return "", err
	}

	return systemProject.Name, nil
}

func (n *ProviderCatalogLauncher) getExternalDNSCatalogID(clusterName string) (string, error) {
	templateVersionID := n.getRangerExternalDNSTemplateID()
	return n.catalogManager.GetSystemAppCatalogID(templateVersionID, clusterName)
}

func (n *ProviderCatalogLauncher) getRangerExternalDNSTemplateID() string {
	return fmt.Sprintf("%s-%s", cutils.SystemLibraryName, templateName)
}

func CopyCreatorID(toAnnotations, fromAnnotations map[string]string) map[string]string {
	if val, exist := fromAnnotations[cattleCreatorIDAnnotationKey]; exist {
		if toAnnotations == nil {
			toAnnotations = make(map[string]string, 2)
		}

		toAnnotations[cattleCreatorIDAnnotationKey] = val
	}
	return toAnnotations
}

func answersDiffer(appAnswers map[string]string, newAnswers map[string]string) bool {
	for key, value := range newAnswers {
		if !strings.EqualFold(appAnswers[key], value) {
			return true
		}
	}
	return false
}

func updateAnswers(appAnswers map[string]string, newAnswers map[string]string) {
	for key, value := range newAnswers {
		appAnswers[key] = value
	}
}

func removeOwnerRef(app *pv3.App, refToRemove metav1.OwnerReference) {
	var newRef []metav1.OwnerReference
	for _, ref := range app.OwnerReferences {
		if ref.Name != refToRemove.Name && ref.UID != refToRemove.UID {
			newRef = append(newRef, ref)
		}
	}
	app.OwnerReferences = newRef
}
