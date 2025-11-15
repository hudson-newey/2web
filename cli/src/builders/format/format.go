package format

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func FormatSolution(args []string) {
	config, err := configs.PrettierConfigLocation()
	pathTargets := builders.EntryTargets(args)

	hasFormatConfig := err == nil
	if hasFormatConfig {
		cliArgs := append(
			[]string{
				"prettier",
				"--write",
				"--config",
				config,
			},
			pathTargets...,
		)

		packages.ExecutePackage(cliArgs...)
	} else {
		cliArgs := append([]string{"prettier", "--write"}, pathTargets...)
		packages.ExecutePackage(cliArgs...)
	}
}
