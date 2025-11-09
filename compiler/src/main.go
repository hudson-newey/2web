package main

import (
	"fmt"
	"hudson-newey/2web/src/builder"
	"hudson-newey/2web/src/builder/cache"
	"hudson-newey/2web/src/cli"
	"os"
	"time"
)

func main() {
	startTime := time.Now()

	cli.ParseArguments()
	isErrorFree := builder.Build()

	if !*cli.GetArgs().IsSilent {
		fmt.Println("\nCompile time:", time.Since(startTime))
	}

	// Ensure that the database connection is properly closed before exiting.
	// I suspect that this would be automatically handled, but I'm not 100% sure.
	//
	// I do not include this in build times because by the time we reach this
	// point, all build assets have been generated and fully usable.
	cache.CloseDBConnection()

	if isErrorFree {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
