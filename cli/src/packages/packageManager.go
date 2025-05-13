package packages

import (
	"os"

	"github.com/hudson-newey/2web-cli/src/cli"
)

type PackageManager = int

const (
	Npm PackageManager = iota
	Pnpm
	Yarn
	Bun
)

func DeterminePackageManager() PackageManager {
	if fileExists("package-lock.json") {
		return Npm
	} else if fileExists("pnpm-lock.yaml") || fileExists("pnpm-lock.yml") {
		return Pnpm
	} else if fileExists("yarn.lock") {
		return Yarn
	} else if fileExists("bun.lock") {
		return Bun
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
