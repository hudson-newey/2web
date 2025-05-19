package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func ServeSolution() {
	viteConfig, err := configs.ViteConfigLocation()

	if err == nil {
		packages.ExecutePackage("vite", "./src/", "--config", viteConfig)
	} else {
		// If there is no vite config, then we want to execute vite without any
		// --config arguments, meaning that Vite should use the default config.
		packages.ExecutePackage("vite", "./src/")
	}
}
