package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/installer"
	"github.com/hudson-newey/2web/_shared/logger"
)

func installCommand(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <package_name>", programName, command)
		logger.PrintError(errorMsg)
	}

	packageName := args[2]
	installer.InstallPackage(packageName)
}
