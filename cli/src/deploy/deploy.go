package deploy

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/deploy/netlify"
	"github.com/hudson-newey/2web/_shared/logger"
)

func DeploySolution(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <location> [arguments]", programName, command)
		logger.PrintError(errorMsg)
	}

	location := args[2]

	switch location {
	case "netlify":
		netlify.Deploy()
	default:
		errorMsg := fmt.Sprintf("Unsupported deployment location: '%s'", location)
		logger.PrintError(errorMsg)
	}
}
