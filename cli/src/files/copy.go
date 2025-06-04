package files

import (
	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/otiai10/copy"
)

func CopyPath(src string, dst string) {
	err := copy.Copy(src, dst)
	if err != nil {
		cli.PrintError(err.Error())
	}

	logCopy(dst)
}
