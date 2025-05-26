package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func TestSolution(args []string) {
	wtrConfig, err := configs.WtrConfigLocation()
	pathTarget := entryTarget(args)

	if err == nil {
		packages.ExecutePackage(
			"wtr",
			"--config",
			wtrConfig,
			pathTarget,
		)
	} else {
		packages.ExecutePackage("wtr", "./src/")
	}
}
