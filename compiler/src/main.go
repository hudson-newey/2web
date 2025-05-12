package main

import (
	"hudson-newey/2web/src/builder"
	"hudson-newey/2web/src/cli"
)

func main() {
	cli.ParseArguments()
	builder.Build()
}
