package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/db"
	"github.com/hudson-newey/2web/_shared/logger"
)

func dbCommand(programName string, command string, args []string) {
	if len(args) < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <sub_command> [arguments]", programName, command)
		logger.PrintError(errorMsg)
		return
	}

	subCommand := args[2]

	switch subCommand {
	case "init":
		db.InitDatabase()
	case "migrate", "m":
		db.RunMigration()
	default:
		errorMsg := fmt.Sprintf("Unknown sub command: '%s'", subCommand)
		logger.PrintError(errorMsg)
	}
}
