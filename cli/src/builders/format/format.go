package format

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func FormatSolution(args []string) {
	config, err := configs.OxFmtConfigLocation()
	pathTargets := builders.EntryTargets(args)

	hasFormatConfig := err == nil
	if hasFormatConfig {
		cliArgs := append([]string{"oxfmt", "-c", config}, pathTargets...)

		packages.ExecutePackage(cliArgs...)
	} else {
		cliArgs := append([]string{"oxfmt"}, pathTargets...)
		packages.ExecutePackage(cliArgs...)
	}
}
