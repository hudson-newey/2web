package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func BuildSolution() {
	viteConfig, err := configs.ViteConfigLocation()

	if err == nil {
		packages.ExecutePackage("vite", "build", "./src/", "--config", viteConfig)
	} else {
		// If there is no Vite config, in the current project, we call Vite without
		// the --config parameter, meaning that it should use the default Vite
		// config.
		packages.ExecutePackage("vite", "build", "./src/")
	}
}
