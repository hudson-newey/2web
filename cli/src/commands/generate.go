package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/templates"
)

func Generate(programName string, command string, args []string) {
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

	if template == "component" || template == "c" {
		templates.ComponentTemplate(templateName)
	} else if template == "service" || template == "s" {
		templates.ServiceTemplate(templateName)
	} else if template == "model" || template == "m" {
		templates.ModelTemplate(templateName)
	} else if template == "aspect" || template == "a" {
		templates.AspectTemplate(templateName)
	} else if template == "interceptor" || template == "i" {
		templates.InterceptorTemplate(templateName)
	} else if template == "page" || template == "p" {
		templates.PageTemplate(templateName)
	} else {
		errorMsg := fmt.Sprintf("unrecognized generate template: '%s'", template)
		cli.PrintError(1, errorMsg)
	}
}
