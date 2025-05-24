package main

import (
	"fmt"
	"hudson-newey/2web/src/builder"
	"hudson-newey/2web/src/cli"
	"time"
)

func main() {
	startTime := time.Now()

	cli.ParseArguments()
	builder.Build()

	if !*cli.GetArgs().IsSilent {
		fmt.Println("\nCompile time:", time.Since(startTime))
	}
}
