package packages

import "github.com/hudson-newey/2web-cli/src/shell"

// Allows you to run a package managers "exec" command.
// This function will automatically determine the package manager used in the
// current solution, and use it.
// If no package manager is found or the dependency isn't a dependency of the
// current solution, the "npm dlx" command will be run instead.
func ExecutePackage(args ...string) {
	packageManager := DeterminePackageManager()

	binaryPath := ""
	switch packageManager {
	case Npm:
		binaryPath = "npm"
	case Pnpm:
		binaryPath = "pnpm"
	case Yarn:
		binaryPath = "yarn"
	case Bun:
		binaryPath = "bun"
	}

	shellCommand := []string{binaryPath, "exec"}
	shellCommand = append(shellCommand, args...)

	shell.ExecuteCommand(shellCommand...)
}
