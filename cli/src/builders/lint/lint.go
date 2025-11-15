package lint

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func LintSolution(args []string) {
	oxLintConfig, err := configs.OxLintConfigLocation()
	pathTargets := builders.EntryTargets(args)

	hasOxLintConfig := err == nil
	if hasOxLintConfig {
		args := append(
			[]string{
				"oxlint",
				"lint",
				"--config",
				oxLintConfig,
			},
			pathTargets...,
		)

		packages.ExecutePackage(args...)
	} else {
		args := append([]string{"oxlint", "lint"}, pathTargets...)
		packages.ExecutePackage(args...)
	}
}
