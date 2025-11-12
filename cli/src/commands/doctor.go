package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/doctor"
)

func doctorCommand(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <sub_command>", programName, command)
		cli.PrintError(errorMsg)
	}

	subCommand := args[2]
	switch subCommand {
	case "check", "c":
		doctor.RunDoctor()
	case "check-dependencies", "cd":
		doctor.CheckDependencies()
	}
}
