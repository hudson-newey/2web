package build

import "github.com/hudson-newey/2web/_shared/shell"

func buildWithInbuiltCompiler(
	compilerPath string,
	inPath string,
	outPath string,
) {
	shell.ExecuteSync(compilerPath, "-i", inPath, "-o", outPath)
}
