package build

import (
	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/shell"
)

func copyAssetsOnly(inPath string, outPath string) {
	cli.PrintWarning("2webc compiler not found. Copying assets without compilation.")

	// For some reason copying files is not a trivial task in Go.
	// TODO: Replace this shell script hack with a proper Go implementation.
	shell.ExecuteCommand("cp", "-r", inPath+"/*", outPath)
}
