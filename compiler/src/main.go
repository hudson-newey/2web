package main

import (
	"hudson-newey/2web/src/actions"
	"hudson-newey/2web/src/builder"
	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/cli"
	"os"
)

func main() {
	cli.ParseArguments()

	if cli.GetArgs().Listen {
		// Once you init the action server, the program will enter an infinite loop.
		actions.InitActionServer()
	}

	// Only build if the listen flag is not being used.
	isErrorFree := builder.Build()

	// Ensure that the database connection is properly closed before exiting.
	// I suspect that this would be automatically handled, but I'm not 100% sure.
	//
	// I do not include this in build times because by the time we reach this
	// point, all build assets have been generated and fully usable.
	//
	// Additionally, we want to keep the DB connection open if the action server
	// is running.
	cache.CloseDBConnection()

	if isErrorFree {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
