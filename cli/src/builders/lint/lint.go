package lint

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func LintSolution(args []string) {
	config, err := configs.OxLintConfigLocation()
	pathTargets := builders.EntryTargets(args)

	hasLintConfig := err == nil
	if hasLintConfig {
		cliArgs := append(
			[]string{
				"oxlint",
				"lint",
				"--config",
				config,
			},
			pathTargets...,
		)

		packages.ExecutePackage(cliArgs...)
	} else {
		cliArgs := append([]string{"oxlint", "lint"}, pathTargets...)
		packages.ExecutePackage(cliArgs...)
	}
}
