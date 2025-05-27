package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/templates"
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
		templates.ComponentTemplate(templateName)
	case "service", "s":
		templates.ServiceTemplate(templateName)
	case "model", "m":
		templates.ModelTemplate(templateName)
	case "aspect", "a":
		templates.AspectTemplate(templateName)
	case "interceptor", "i":
		templates.InterceptorTemplate(templateName)
	case "page", "p":
		templates.PageTemplate(templateName)
	// Templates below this point require @two-web/kit
	case "guard", "g":
		templates.GuardTemplate(templateName)
	default:
		cli.PrintError(1, fmt.Sprintf("unrecognized generate template: '%s'", template))
	}
}
