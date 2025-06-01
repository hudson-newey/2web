package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
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
