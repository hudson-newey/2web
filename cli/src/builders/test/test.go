package test

import (
	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func TestSolution(args []string) {
	wtrConfig, err := configs.WtrConfigLocation()
	pathTargets := builders.EntryTargets(args)

	if err == nil {
		args := append([]string{"wtr", "--config", wtrConfig}, pathTargets...)
		packages.ExecuteWithoutFallback(args...)
	} else {
		args := append([]string{"wtr"}, pathTargets...)
		packages.ExecuteWithoutFallback(args...)
	}
}
