package files

import (
	"github.com/hudson-newey/2web/_shared/logger"
	"github.com/otiai10/copy"
)

func CopyPath(src string, dst string) {
	err := copy.Copy(src, dst)
	if err != nil {
		logger.PrintError(err.Error())
	}

	logCopy(dst)
}
