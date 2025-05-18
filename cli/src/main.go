package main

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/commands"
	"github.com/hudson-newey/2web-cli/src/packages"
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

	if command == "help" || command == "--help" {
		commands.PrintHelpDocs()
		return
	}

	if command == "version" || command == "--version" {
		commands.PrintVersion()
		return
	}

	// all commands below this point must be run in an npm project
	if !packages.HasPackageJson() {
		cli.PrintError(2, "the 2web cli must be run in a directory with a package.json")
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

	if command == "lint" {
		commands.Lint()
		return
	}

	errorMsg := fmt.Sprintf("unrecognized command: '%s'", command)
	cli.PrintError(1, errorMsg)
}
