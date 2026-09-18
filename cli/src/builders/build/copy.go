package build

import (
	"github.com/hudson-newey/2web/_shared/logger"
	"github.com/hudson-newey/2web/_shared/shell"
)

func copyAssetsOnly(inPath string, outPath string) {
	logger.PrintWarning("2webc compiler not found. Copying assets without compilation.")

	// For some reason copying files is not a trivial task in Go.
	// TODO: Replace this shell script hack with a proper Go implementation.
	shell.ExecuteSync("cp", "-r", inPath+"/*", outPath)
}
