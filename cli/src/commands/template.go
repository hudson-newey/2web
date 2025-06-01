package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
)

func template(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <template>", programName, command)
		cli.PrintError(1, errorMsg)
	}

	template := args[2]

	switch template {
	case "ssr":
	}
}
