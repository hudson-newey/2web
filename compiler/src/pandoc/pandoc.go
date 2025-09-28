package pandoc

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"os/exec"
)

func ConvertFormat(
	content string,
	fromFormat string,
	toFormat string,
) (string, error) {
	if !isPandocInstalled() {
		errorMsg := fmt.Sprintf(
			"Pandoc is required to compile '%s' formats is not installed. "+
				"Please install pandoc https://pandoc.org/installing.html",
			fromFormat,
		)
		cli.HardError(errorMsg)
		return "", nil
	}

	return content, nil
}

func isPandocInstalled() bool {
	pandocPath, err := exec.LookPath("pandoc")
	return err == nil && pandocPath != ""
}
