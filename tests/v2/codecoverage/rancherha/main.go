package main

import (
	"github.com/ranger/ranger/tests/framework/clients/corral"
	"github.com/ranger/ranger/tests/framework/pkg/session"
	"github.com/sirupsen/logrus"
)

func main() {
	testSession := session.NewSession()

	corralConfig := corral.CorralConfigurations()
	err := corral.SetupCorralConfig(corralConfig.CorralConfigVars, corralConfig.CorralConfigUser, corralConfig.CorralSSHPath)
	if err != nil {
		logrus.Fatalf("error setting up corral: %v", err)
	}
	configPackage := corral.CorralPackagesConfig()

	path := configPackage.CorralPackageImages["rangertestcoverage"]
	_, err = corral.CreateCorral(testSession, "rangertestcoverage", path, true, configPackage.HasCleanup)
	if err != nil {
		logrus.Errorf("error creating corral: %v", err)
	}
}
