package install

import (
	"os"

	"github.com/hudson-newey/2web-cli/src/cli"
)

type packageManager = int

const (
	npm packageManager = iota
	pnpm
	yarn
	bun
)

func InstallPackage(packageName string) {
	installationMethod := determinePackageManager()

	switch installationMethod {
	case npm:
		installNpmPackage(packageName)
	case pnpm:
		installPnpmPackage(packageName)
	case yarn:
		installYarnPackage(packageName)
	case bun:
		installBunPackage(packageName)
	}
}

func determinePackageManager() packageManager {
	if fileExists("package-lock.json") {
		return npm
	} else if fileExists("pnpm-lock.yaml") || fileExists("pnpm-lock.yml") {
		return pnpm
	} else if fileExists("yarn.lock") {
		return yarn
	} else if fileExists("bun.lock") {
		return bun
	}

	cli.PrintError(2, "could not determine package manager")
	panic("")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil // true if no error (file exists)
}
