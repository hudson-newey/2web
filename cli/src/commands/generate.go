package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/generators"
)

func generate(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <generator> <name>", programName, command)
		cli.PrintError(errorMsg)
	}

	generator := args[2]
	if argsLen < 4 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s %s <name>", programName, command, generator)
		cli.PrintError(errorMsg)
	}

	templateName := args[3]

	switch generator {
	case "component", "c":
		generators.ComponentGenerator(templateName)
	case "service", "s":
		generators.ServiceGenerator(templateName)
	case "aspect", "a":
		generators.AspectGenerator(templateName)
	case "interceptor", "i":
		generators.InterceptorGenerator(templateName)
	case "page", "p":
		generators.PageGenerator(templateName)
	// generators below this point require @two-web/kit
	case "guard", "g":
		generators.GuardGenerator(templateName)
	case "model", "m":
		generators.ModelGenerator(templateName)
	case "enum", "e":
		generators.EnumGenerator(templateName)
	default:
		cli.PrintError(fmt.Sprintf("unrecognized generate template: '%s'", generator))
	}
}
