package packages

import "github.com/hudson-newey/2web-cli/src/shell"

// Allows you to run a package managers "exec" command.
// This function will automatically determine the package manager used in the
// current solution, and use it.
// If no package manager is found or the dependency isn't a dependency of the
// current solution, the "npm dlx" command will be run instead.
//
// Even though a package manager might change the command syntax to execute
// installed package binaries, at the time of writing, all of the supported
// package managers use the same format.
// I have therefore reasoned that it is likely that the package manager
// maintainers will not attempt to break this convention, and if they do end
// up changing the format, they will provide enough time for me to adjust
// the shell command.
func ExecutePackage(args ...string) {
	packageManager := DeterminePackageManager()

	if packageManager == None {
		executeNpx(args...)
		return
	}

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

// If no local package manager is installed, we can run most commands using the
// "npm dlx" command.
// This will temporarily download the package to run the command.
// This is not recommended because it can result in unpredictable environments.
//
// Although, supporting this use case makes the 2web cli useful for mocking up
// projects outside of the 2web framework, making the cli more framework
// agnostic.
func executeNpx(args ...string) {
	shellCommand := []string{"npx"}
	shellCommand = append(shellCommand, args...)

	shell.ExecuteCommand(shellCommand...)
}
