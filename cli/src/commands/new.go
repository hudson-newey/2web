package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/templates"
)

func new(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <project_name>", programName, command)
		cli.PrintError(errorMsg)
	}

	projectName := args[2]
	templates.NewTemplate(projectName)
}
