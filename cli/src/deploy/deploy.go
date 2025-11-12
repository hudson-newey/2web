package deploy

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/deploy/netlify"
)

func DeploySolution(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <location> [arguments]", programName, command)
		cli.PrintError(errorMsg)
	}

	location := args[2]

	switch location {
	case "netlify":
		netlify.Deploy()
	default:
		errorMsg := fmt.Sprintf("Unsupported deployment location: '%s'", location)
		cli.PrintError(errorMsg)
	}
}
