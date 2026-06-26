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
	Deno
	None
)

func DeterminePackageManager() PackageManager {
	if fileExists("bun.lock") {
		return Bun
	} else if fileExists("deno.lock") {
		return Deno
	} else if fileExists("package-lock.json") {
		return Npm
	} else if fileExists("pnpm-lock.yaml") || fileExists("pnpm-lock.yml") {
		return Pnpm
	} else if fileExists("yarn.lock") || fileExists(".pnp.js") || fileExists(".pnp.cjs") {
		return Yarn
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
