package devserver

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	pollTime := getPollTime(args)

	server.Run(inPath, outPath, pollTime)
}

// getPollTime extracts the --pollTime flag value from args.
// Returns the default value of 100ms if not specified or invalid.
// Valid range is 1-10000ms. Values outside this range fall back to default.
func getPollTime(args []string) int {
	const (
		defaultPollTime = 100
		minPollTime     = 1
		maxPollTime     = 10000
	)

	for i, arg := range args {
		if strings.HasPrefix(arg, "--pollTime=") {
			valueStr := strings.TrimPrefix(arg, "--pollTime=")
			if value, err := strconv.Atoi(valueStr); err == nil && value >= minPollTime && value <= maxPollTime {
				return value
			}
		} else if arg == "--pollTime" && i+1 < len(args) {
			if value, err := strconv.Atoi(args[i+1]); err == nil && value >= minPollTime && value <= maxPollTime {
				return value
			}
		}
	}
	return defaultPollTime
}

func serveSsr() {
	shell.ExecuteCommand("node", "./server/ssr.ts")
}
