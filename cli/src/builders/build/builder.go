package build

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/ssr"
)

func BuildSolution(args []string) {
	viteConfig, err := configs.ViteConfigLocation()
	pathTarget := builders.EntryTarget(args)

	if err == nil {
		if ssr.HasSsrTarget() {
			packages.ExecutePackage("vite", "build", pathTarget, "./server/", "--config", viteConfig)
		} else {
			packages.ExecutePackage("vite", "build", pathTarget, "--config", viteConfig)
		}
	} else {
		// If there is no Vite config, in the current project, we call Vite without
		// the --config parameter, meaning that it should use the default Vite
		// config.
		if ssr.HasSsrTarget() {
			packages.ExecutePackage("vite", "build", pathTarget, "./server/")
		} else {
			packages.ExecutePackage("vite", "build", pathTarget)
		}
	}
}
