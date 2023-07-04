package main

import (
	"github.com/ranger/ranger/tests/framework/clients/corral"
	"github.com/ranger/ranger/tests/framework/clients/ranger"
	"github.com/ranger/ranger/tests/framework/extensions/pipeline"
	"github.com/ranger/ranger/tests/framework/pkg/config"
	"github.com/ranger/ranger/tests/framework/pkg/environmentflag"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/ranger/ranger/tests/v2/validation/pipeline/rangerha/corralha"
	"github.com/sirupsen/logrus"
)

func main() {
	corralRangerHA := new(corralha.CorralRangerHA)
	config.LoadConfig(corralha.CorralRangerHAConfigConfigurationFileKey, corralRangerHA)

	corralSession := session.NewSession()

	corralConfig := corral.CorralConfigurations()
	err := corral.SetupCorralConfig(corralConfig.CorralConfigVars, corralConfig.CorralConfigUser, corralConfig.CorralSSHPath)
	if err != nil {
		logrus.Fatalf("error setting up corral: %v", err)
	}

	configPackage := corral.CorralPackagesConfig()

	environmentFlags := environmentflag.NewEnvironmentFlags()
	environmentflag.LoadEnvironmentFlags(environmentflag.ConfigurationFileKey, environmentFlags)
	installRanger := environmentFlags.GetValue(environmentflag.InstallRanger)

	logrus.Infof("installRanger value is %t", installRanger)

	if installRanger {
		path := configPackage.CorralPackageImages[corralRangerHA.Name]
		corralName := corralRangerHA.Name

		_, err = corral.CreateCorral(corralSession, corralName, path, true, configPackage.HasCleanup)
		if err != nil {
			logrus.Errorf("error creating corral: %v", err)
		}

		bootstrapPassword, err := corral.GetCorralEnvVar(corralName, "bootstrap_password")
		if err != nil {
			logrus.Errorf("error getting the bootstrap password: %v", err)
		}

		if configPackage.HasSetCorralSSHKeys {
			privateKey, err := corral.GetCorralEnvVar(corralName, "corral_private_key")
			if err != nil {
				logrus.Errorf("error getting the corral's private key: %v", err)
			}
			logrus.Infof("Corral Private Key: %s", privateKey)

			publicKey, err := corral.GetCorralEnvVar(corralName, "corral_public_key")
			if err != nil {
				logrus.Errorf("error getting the corral's public key: %v", err)
			}
			logrus.Infof("Corral Public Key: %s", publicKey)

			err = corral.UpdateCorralConfig("corral_private_key", privateKey)
			if err != nil {
				logrus.Errorf("error setting the corral's private key: %v", err)
			}

			err = corral.UpdateCorralConfig("corral_public_key", publicKey)
			if err != nil {
				logrus.Errorf("error setting the corral's public key: %v", err)
			}
		}

		rangerConfig := new(ranger.Config)
		config.LoadConfig(ranger.ConfigurationFileKey, rangerConfig)

		token, err := pipeline.CreateAdminToken(bootstrapPassword, rangerConfig)
		if err != nil {
			logrus.Errorf("error creating the admin token: %v", err)
		}
		rangerConfig.AdminToken = token
		config.UpdateConfig(ranger.ConfigurationFileKey, rangerConfig)
		rangerSession := session.NewSession()
		client, err := ranger.NewClient(rangerConfig.AdminToken, rangerSession)
		if err != nil {
			logrus.Errorf("error creating the ranger client: %v", err)
		}

		err = pipeline.PostRangerInstall(client, rangerConfig.AdminPassword)
		if err != nil {
			logrus.Errorf("error during post ranger install: %v", err)
		}
	} else {
		logrus.Infof("Skipped Ranger Install because installRanger is %t", installRanger)
	}
}
