package main

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/commands"
)

func main() {
	argArray := os.Args
	programName := argArray[0]
	argsLen := len(argArray)

	if argsLen < 2 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s <command> [arguments]", programName)
		cli.PrintError(1, errorMsg)
	}

	command := os.Args[1]

	if command == "new" || command == "n" {
		commands.New(programName, command, argArray)
		return
	}

	if command == "generate" || command == "g" {
		commands.Generate(programName, command, argArray)
		return
	}

	if command == "install" || command == "i" {
		commands.Install(programName, command, argArray)
		return
	}

	if command == "deploy" {
		commands.Deploy()
		return
	}

	if command == "serve" {
		commands.Serve()
		return
	}

	if command == "build" {
		commands.Build()
		return
	}

	errorMsg := fmt.Sprintf("unrecognized command: '%s'", command)
	cli.PrintError(1, errorMsg)
}
