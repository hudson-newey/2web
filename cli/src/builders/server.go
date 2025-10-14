package builders

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/shell"
)

func ServeSolution(args []string) {
	if hasSsrTarget() {
		serveSsr()
	} else {
		serveBrowser(args)
	}
}

func serveBrowser(args []string) {
	viteConfig, err := configs.ViteConfigLocation()
	pathTarget := entryTarget(args)

	// Check that the path target actually exists.
	// If it does not, we want to log a warning.
	if _, err := os.Stat(pathTarget); os.IsNotExist(err) {
		warningMsg := fmt.Sprintf("the specified path target does not exist: '%s'", pathTarget)
		cli.PrintWarning(warningMsg)
	}

	if err == nil {
		packages.ExecutePackage("vite", pathTarget, "--config", viteConfig)
	} else {
		// If there is no vite config, then we want to execute vite without any
		// --config arguments, meaning that Vite should use the default config.
		packages.ExecutePackage("vite", pathTarget)
	}
}

func serveSsr() {
	shell.ExecuteCommand("node", "./server/ssr.ts")
}
