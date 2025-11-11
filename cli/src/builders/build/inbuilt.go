package build

import "github.com/hudson-newey/2web-cli/src/shell"

func buildWithInbuiltCompiler(
	compilerPath string,
	inPath string,
	outPath string,
) {
	// TODO: Remove this --no-cache flag once caching is stable
	shell.ExecuteCommand(compilerPath, "-i", inPath, "-o", outPath, "--no-cache")
}
