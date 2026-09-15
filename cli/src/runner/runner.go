package runner

import (
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/shell"
	"github.com/hudson-newey/2web/_shared/logger"
)

// Executes a JavaScript / TypeScript file
func ExecuteScript(filePath string, args ...string) {
	packageManager := packages.DeterminePackageManager()

	if packageManager == packages.None {
		logger.PrintError("could not find a package manager in the current solution.")
		return
	}

	shellCommand := runtimeCommand(packageManager)
	shellCommand = append(shellCommand, filePath)
	shellCommand = append(shellCommand, args...)

	shell.ExecuteCommand(shellCommand...)
}

// Returns the base commands needed to execute a file using the projects
// detected runtime.
// The returned array should be owned by the consumer and not be used in any
// other state.
func runtimeCommand(pm packages.PackageManager) []string {
	switch pm {
	case packages.Bun:
		return []string{"bun"}
	case packages.Deno:
		return []string{"deno", "-A"}
	case packages.Npm, packages.Pnpm, packages.Yarn:
		// Node can execute TypeScript directly from node 22.6 (with the
		// type stripping flag) and 23.6+ (natively). Older node versions
		// reject the flag, but the flag is a no-op on modern runtimes.
		return []string{"node", "--experimental-strip-types"}
	}

	logger.PrintError("Could not determine runtime")
	panic("unreachable")
}
