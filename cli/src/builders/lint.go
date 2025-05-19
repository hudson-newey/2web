package builders

import (
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func LintSolution() {
	packages.ExecutePackage(
		"eslint",
		"--config",
		configs.EslintConfigLocation(),
		"./src/",
	)
}
