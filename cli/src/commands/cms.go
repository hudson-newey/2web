package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/cms"
)

func cmsCommand(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <sub_command> [arguments]", programName, command)
		cli.PrintError(errorMsg)
		return
	}

	subCommand := args[2]
	switch subCommand {
	case "add", "a":
		cms.AddCmsSource(args)
	case "view", "v":
		cms.ViewCmsSource()
	case "sync", "s":
		cms.SyncCmsSources(args)
	case "remove", "rm":
		cms.RemoveCmsSource()
	default:
		errorMsg := fmt.Sprintf("Unknown cms command: '%s'", subCommand)
		cli.PrintError(errorMsg)
	}
}
