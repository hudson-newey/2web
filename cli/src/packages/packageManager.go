package packages

import (
	"os"
)

type PackageManager = int

const (
	Npm PackageManager = iota
	Pnpm
	Yarn
	Bun
	None
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

	return None
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil // true if no error (file exists)
}
