package data

import (
	"github.com/ranger/ranger/pkg/auth/providers/activedirectory"
	"github.com/ranger/ranger/pkg/auth/providers/azure"
	"github.com/ranger/ranger/pkg/auth/providers/github"
	"github.com/ranger/ranger/pkg/auth/providers/googleoauth"
	"github.com/ranger/ranger/pkg/auth/providers/keycloakoidc"
	"github.com/ranger/ranger/pkg/auth/providers/ldap"
	localprovider "github.com/ranger/ranger/pkg/auth/providers/local"
	"github.com/ranger/ranger/pkg/auth/providers/oidc"
	"github.com/ranger/ranger/pkg/auth/providers/saml"
	client "github.com/ranger/ranger/pkg/client/generated/management/v3"
	"github.com/ranger/ranger/pkg/controllers/management/auth"
	v3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
	"github.com/ranger/ranger/pkg/types/config"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AuthConfigs(management *config.ManagementContext) error {
	if err := addAuthConfig(github.Name, client.GithubConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(activedirectory.Name, client.ActiveDirectoryConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(azure.Name, client.AzureADConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(ldap.OpenLdapName, client.OpenLdapConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(ldap.FreeIpaName, client.FreeIpaConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(saml.PingName, client.PingConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(saml.ADFSName, client.ADFSConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(saml.KeyCloakName, client.KeyCloakConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(saml.OKTAName, client.OKTAConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(saml.ShibbolethName, client.ShibbolethConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(googleoauth.Name, client.GoogleOauthConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(oidc.Name, client.OIDCConfigType, false, management); err != nil {
		return err
	}

	if err := addAuthConfig(keycloakoidc.Name, client.KeyCloakOIDCConfigType, false, management); err != nil {
		return err
	}

	return addAuthConfig(localprovider.Name, client.LocalConfigType, true, management)
}

func addAuthConfig(name, aType string, enabled bool, management *config.ManagementContext) error {
	annotations := make(map[string]string)
	if name == azure.Name {
		annotations[azure.GraphEndpointMigratedAnnotation] = "true"
	}
	annotations[auth.CleanupAnnotation] = auth.CleanupRangerLocked

	_, err := management.Management.AuthConfigs("").ObjectClient().Create(&v3.AuthConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
		},
		Type:    aType,
		Enabled: enabled,
	})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}
