package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func FormatSolution(args []string) {
	biomeConfig, err := configs.BiomeConfigLocation()
	pathTarget := entryTarget(args)

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
		if hasSsrTarget() {
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
