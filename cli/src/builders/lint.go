package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func LintSolution() {
	eslintConfig, err := configs.EslintConfigLocation()

	if err == nil {
		packages.ExecutePackage(
			"eslint",
			"--config",
			eslintConfig,
			"./src/",
		)
	} else {
		packages.ExecutePackage("eslint", "./src/")
	}
}
