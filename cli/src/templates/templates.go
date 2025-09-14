package templates

import (
	"errors"
	"os"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/files"
)

const templatesPath string = "./node_modules/@two-web/cli/templates/"

func copyFromTemplates(name string, dst string) {
	if !areTemplatesAvailable() {
		cli.PrintError("unable to apply template. \"@two-web/cli/templates\" is not available.")
	}

	path := templatesPath + name
	files.CopyPath(path, dst)
}

func areTemplatesAvailable() bool {
	_, err := os.Stat(templatesPath)
	return !errors.Is(err, os.ErrNotExist)
}
