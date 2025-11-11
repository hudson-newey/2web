package commands

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/builders/build"
	devserver "github.com/hudson-newey/2web-cli/src/builders/devServer"
	"github.com/hudson-newey/2web-cli/src/builders/format"
	"github.com/hudson-newey/2web-cli/src/builders/lint"
	"github.com/hudson-newey/2web-cli/src/builders/test"
	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/deploy"
	"github.com/hudson-newey/2web-cli/src/packages"
)

func ProcessInvocation(args []string) {
	programName := args[0]
	argsLen := len(args)

	if argsLen < 2 {
		errorMsg := fmt.Sprintf("invalid arguments:\n\texpected: %s <command> [arguments]", programName)
		cli.PrintError(errorMsg)
	}

	command := os.Args[1]

	if command == "new" || command == "n" {
		new(programName, command, args)
		return
	}

	if command == "help" || command == "--help" {
		printHelpDocs()
		return
	}

	if command == "version" || command == "--version" {
		printVersion()
		return
	}

	// all commands below this point must be run in an npm project
	if !packages.HasPackageJson() {
		cli.PrintWarning("you are running the 2web cli in a directory that does not have a package.json. This may result in unpredictable state.")
	}

	if command == "generate" || command == "g" {
		generate(programName, command, args)
		return
	}

	if command == "template" || command == "t" {
		template(programName, command, args)
		return
	}

	if command == "install" || command == "i" {
		install(programName, command, args)
		return
	}

	if command == "doctor" {
		doctorCommand(programName, command, args)
		return
	}

	// We pass the commands into the "serve" builder because we support specifying
	// a path to serve as the third optional argument.
	// If not specified, the "./src/" directory will be the target.
	if command == "serve" || command == "s" {
		devserver.ServeSolution(args)
		return
	}

	if command == "build" {
		build.BuildSolution(args)
		return
	}

	if command == "lint" {
		lint.LintSolution(args)
		return
	}

	if command == "format" {
		format.FormatSolution(args)
		return
	}

	if command == "test" {
		test.TestSolution(args)
	}

	if command == "deploy" {
		deploy.DeploySolution()
		return
	}

	errorMsg := fmt.Sprintf("unrecognized command: '%s'", command)
	cli.PrintError(errorMsg)
}
