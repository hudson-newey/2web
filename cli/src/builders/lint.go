package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func LintSolution(args []string) {
	eslintConfig, err := configs.EslintConfigLocation()
	pathTarget := entryTarget(args)

	if err == nil {
		packages.ExecutePackage(
			"eslint",
			"--config",
			eslintConfig,
			pathTarget,
		)
	} else {
		packages.ExecutePackage("eslint", "./src/")
	}
}
