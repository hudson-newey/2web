package test

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func TestSolution(args []string) {
	wtrConfig, err := configs.WtrConfigLocation()
	pathTarget := builders.EntryTarget(args)

	if err == nil {
		packages.ExecuteWithoutFallback(
			"wtr",
			"--config",
			wtrConfig,
			pathTarget,
		)
	} else {
		packages.ExecuteWithoutFallback("wtr", "./src/")
	}
}
