package lint

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/ssr"
)

func LintSolution(args []string) {
	oxLintConfig, err := configs.OxLintConfigLocation()
	pathTarget := builders.EntryTarget(args)

	if err == nil {
		packages.ExecutePackage(
			"oxlint",
			"lint",
			"--config",
			oxLintConfig,
			pathTarget,
		)
	} else {
		if ssr.HasSsrTarget() {
			packages.ExecutePackage(
				"oxlint",
				"lint",
				"--config",
				oxLintConfig,
				"./src/",
				"./server/",
			)
		} else {
			packages.ExecutePackage(
				"oxlint",
				"lint",
				"--config",
				oxLintConfig,
				"./src/",
			)
		}
	}
}
