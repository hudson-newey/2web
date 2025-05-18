package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func BuildSolution() {
	// Even though a package manager might change the command syntax to execute
	// installed package binaries, at the time of writing, all of the supported
	// package managers use the same format.
	// I have therefore reasoned that it is likely that the package manager
	// maintainers will not attempt to break this convention, and if they do end
	// up changing the format, they will provide enough time for me to adjust
	// the shell command.
	packages.ExecutePackage(
		"vite",
		"build",
		"--config",
		configs.ViteConfigLocation(),
	)
}
