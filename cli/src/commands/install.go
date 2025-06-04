package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/installer"
)

func install(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <package_name>", programName, command)
		cli.PrintError(errorMsg)
	}

	packageName := args[2]
	installer.InstallPackage(packageName)
}
