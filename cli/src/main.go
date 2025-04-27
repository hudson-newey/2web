package main

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/templates"
)

func main() {
	programName := os.Args[0]
	argsLen := len(os.Args)

	if argsLen < 2 {
		errorMsg := fmt.Errorf("invalid arguments:\n\texpected: %s <command> [arguments]", programName)
		panic(errorMsg)
	}

	command := os.Args[1]

	if command == "new" || command == "n" {
		if argsLen < 3 {
			errorMsg := fmt.Errorf("invalid arguments:\n\texpected: %s %s <project_name>", programName, command)
			panic(errorMsg)
		}

		projectName := os.Args[2]
		templates.NewTemplate(projectName)

		return
	}

	if command == "generate" || command == "g" {
		if argsLen < 3 {
			errorMsg := fmt.Errorf("invalid arguments:\n\texpected: %s %s <template> <template_name>", programName, command)
			panic(errorMsg)
		}

		template := os.Args[2]
		if argsLen < 4 {
			errorMsg := fmt.Errorf("invalid arguments:\n\texpected: %s %s %s <template_name>", programName, command, template)
			panic(errorMsg)
		}

		templateName := os.Args[3]

		if template == "component" || template == "c" {
			templates.ComponentTemplate(templateName)
		} else if template == "service" || template == "s" {
			templates.ServiceTemplate(templateName)
		}
	}
}
