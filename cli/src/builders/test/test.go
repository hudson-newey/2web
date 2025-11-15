package test

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func TestSolution(args []string) {
	pathTargets := builders.EntryTargets(args)

	cliArgs := append([]string{"playwright", "test"}, pathTargets...)
	packages.ExecutePackage(cliArgs...)
}
