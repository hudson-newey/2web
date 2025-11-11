package build

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/ssr"
)

func buildWithVite(inPath string) {
	viteConfig, err := configs.ViteConfigLocation()
	if err != nil {
		// This should only panic if the file is deleted after we checked for its
		// existence.
		panic(err)
	}

	if ssr.HasSsrTarget() {
		packages.ExecutePackage("vite", "build", inPath, "./server/", "--config", viteConfig)
	} else {
		packages.ExecutePackage("vite", "build", inPath, "--config", viteConfig)
	}
}
