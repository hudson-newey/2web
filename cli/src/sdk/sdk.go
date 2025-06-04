package sdk

import (
	"errors"
	"os"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/files"
)

const sdkPath string = "./node_modules/@two-web/sdk/"

func CopyFromSdk(name string, dst string) {
	if !isSdkInstalled() {
		cli.PrintError("unable to apply template. \"@two-web/sdk\" is not installed.")
	}

	path := sdkPath + name
	files.CopyPath(path, dst)
}

func isSdkInstalled() bool {
	_, err := os.Stat(sdkPath)
	return !errors.Is(err, os.ErrNotExist)
}
