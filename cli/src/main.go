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

	if command == "component" || command == "c" {
		if argsLen < 3 {
			errorMsg := fmt.Errorf("invalid arguments:\n\texpected: %s %s <component_name>", programName, command)
			panic(errorMsg)
		}

		componentName := os.Args[2]
		templates.ComponentTemplate(componentName)
	}
}
