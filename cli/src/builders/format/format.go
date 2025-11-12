package format

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/ssr"
)

func FormatSolution(args []string) {
	prettierConfig, err := configs.PrettierConfigLocation()
	pathTarget := builders.EntryTarget(args)

	if err == nil {
		packages.ExecutePackage(
			"prettier",
			"format",
			"--write",
			"--config",
			prettierConfig,
			pathTarget,
		)
	} else {
		if ssr.HasSsrTarget() {
			packages.ExecutePackage(
				"prettier",
				"format",
				"--write",
				"--config",
				prettierConfig,
				"./src/",
				"./server/",
			)
		} else {
			packages.ExecutePackage(
				"prettier",
				"format",
				"--write",
				"--config",
				prettierConfig,
				"./src/",
			)
		}
	}
}
