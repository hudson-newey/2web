package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/generators"
)

func generate(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s <template> <template_name>", programName, command)
		cli.PrintError(1, errorMsg)
	}

	template := args[2]
	if argsLen < 4 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s %s %s <template_name>", programName, command, template)
		cli.PrintError(1, errorMsg)
	}

	templateName := args[3]

	switch template {
	case "component", "c":
		generators.ComponentGenerator(templateName)
	case "service", "s":
		generators.ServiceGenerator(templateName)
	case "model", "m":
		generators.ModelGenerator(templateName)
	case "aspect", "a":
		generators.AspectGenerator(templateName)
	case "interceptor", "i":
		generators.InterceptorGenerator(templateName)
	case "page", "p":
		generators.PageGenerator(templateName)
	// generators below this point require @two-web/kit
	case "guard", "g":
		generators.GuardGenerator(templateName)
	default:
		cli.PrintError(1, fmt.Sprintf("unrecognized generate template: '%s'", template))
	}
}
