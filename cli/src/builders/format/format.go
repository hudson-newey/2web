package format

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/ssr"
)

func FormatSolution(args []string) {
	biomeConfig, err := configs.BiomeConfigLocation()
	pathTarget := builders.EntryTarget(args)

	configPathArg := "--configPath=" + biomeConfig

	if err == nil {
		packages.ExecutePackage(
			"biome",
			"format",
			"--write",
			configPathArg,
			pathTarget,
		)
	} else {
		if ssr.HasSsrTarget() {
			packages.ExecutePackage(
				"biome",
				"format",
				"--write",
				configPathArg,
				"./src/",
				"./server/",
			)
		} else {
			packages.ExecutePackage(
				"biome",
				"format",
				"--write",
				configPathArg,
				"./src/",
			)
		}
	}
}
