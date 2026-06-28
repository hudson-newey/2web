package build

import "github.com/hudson-newey/2web-cli/src/shell"

func buildWithInbuiltCompiler(
	compilerPath string,
	inPath string,
	outPath string,
) {
	shell.ExecuteCommand(compilerPath, "-i", inPath, "-o", outPath)
}
