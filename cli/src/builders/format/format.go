package format

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func FormatSolution(args []string) {
	prettierConfig, err := configs.PrettierConfigLocation()
	pathTargets := builders.EntryTargets(args)

	hasPrettierConfig := err == nil
	if hasPrettierConfig {
		args := append(
			[]string{
				"prettier",
				"--write",
				"--config",
				prettierConfig,
			},
			pathTargets...,
		)

		packages.ExecutePackage(args...)
	} else {
		args := append([]string{"prettier", "--write"}, pathTargets...)
		packages.ExecutePackage(args...)
	}
}
