package main

import (
	"fmt"
	"hudson-newey/2web/src/builder"
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

	if isErrorFree {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
