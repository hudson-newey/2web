package main

import (
	"os"

	"github.com/hudson-newey/2web-cli/src/commands"
)

func main() {
	commands.ProcessInvocation(os.Args)
}
