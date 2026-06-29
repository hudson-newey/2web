package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/templates"
	"github.com/hudson-newey/2web/_shared/logger"
)

func newCommand(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <project_name>", programName, command)
		logger.PrintError(errorMsg)
	}

	projectName := args[2]
	templates.NewTemplate(projectName)
}
