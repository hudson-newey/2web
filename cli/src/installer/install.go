package installer

import (
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web/_shared/logger"
)

func InstallPackage(packageName string) {
	installationMethod := packages.DeterminePackageManager()

	switch installationMethod {
	case packages.Npm:
		installNpmPackage(packageName)
	case packages.Pnpm:
		installPnpmPackage(packageName)
	case packages.Yarn:
		installYarnPackage(packageName)
	case packages.Bun:
		installBunPackage(packageName)
	}

	logger.PrintError("could not determine package manager")
}
