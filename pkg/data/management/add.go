package management

import (
	"context"

	"github.com/ranger/ranger/pkg/auth/data"
	"github.com/ranger/ranger/pkg/types/config"
	"github.com/ranger/ranger/pkg/wrangler"
)

func Add(ctx context.Context, wrangler *wrangler.Context, management *config.ManagementContext) error {
	if err := sshKeyCleanup(management); err != nil {
		return err
	}

	_, err := addRoles(wrangler, management)
	if err != nil {
		return err
	}
	if err := addClusterRoleForNamespacedCRDs(management); err != nil {
		return err
	}

	if err := data.AuthConfigs(management); err != nil {
		return err
	}

	if err := syncCatalogs(management); err != nil {
		return err
	}

	if err := addDefaultPodSecurityPolicyTemplates(management); err != nil {
		return err
	}

	if err = addDefaultPodSecurityAdmissionConfigurationTemplates(management); err != nil {
		return err
	}

	if err := addKontainerDrivers(management); err != nil {
		return err
	}

	if err := addCattleGlobalNamespaces(management); err != nil {
		return err
	}

	return addMachineDrivers(management)
}
