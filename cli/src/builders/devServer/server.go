package devserver

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/server"
	"github.com/hudson-newey/2web-cli/src/shell"
	"github.com/hudson-newey/2web-cli/src/ssr"
)

func ServeSolution(args []string) {
	if ssr.HasSsrTarget() {
		serveSsr()
		return
	}

	if configs.HasViteConfig() {
		serveVite(args)
		return
	}

	serveInbuilt(args)
}

func serveVite(args []string) {
	viteConfig, err := configs.ViteConfigLocation()
	pathTarget := builders.EntryTargets(args)[0]

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

func serveInbuilt(args []string) {
	inPath := builders.EntryTargets(args)[0]
	outPath := builders.OutputTarget(args)

	server.Run(inPath, outPath)
}

func serveSsr() {
	shell.ExecuteCommand("node", "./server/ssr.ts")
}
