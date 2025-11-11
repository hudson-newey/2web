package build

import (
	"os"
	"os/exec"
	"path"

	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/constants"
)

func BuildSolution(args []string) {
	inPath := builders.EntryTarget(args)
	outPath := builders.OutputTarget(args)

	BuildPath(inPath, outPath)
}

func BuildPath(inPath string, outPath string) {
	if configs.HasViteConfig() {
		buildWithVite(inPath)
		return
	}

	compilerPath, err := twoWebCompilerPath()
	if err == nil {
		buildWithInbuiltCompiler(compilerPath, inPath, outPath)
		return
	}

	copyAssetsOnly(inPath, outPath)
}

func twoWebCompilerPath() (string, error) {
	exeName := "2webc"

	localPath := path.Join(constants.LocalInstallPath, exeName)
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	// Search in the PATH environment variable for a globally installed 2webc
	// binary.
	globalPath, err := exec.LookPath(exeName)
	if err == nil {
		return globalPath, nil
	}

	return "", os.ErrNotExist
}
