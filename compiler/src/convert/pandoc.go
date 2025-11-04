package convert

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"os/exec"
)

var emptyFile = []byte{}

// Uses Pandoc to convert between markup formats.
// Pandoc must be installed on the system for this to work.
func ConvertFormat(
	content []byte,
	fromFormat string,
	toFormat string,
) ([]byte, error) {
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		errorMsg := fmt.Sprintf(
			"Pandoc is required to compile '%s' formats, "+
				"but could not be found in the PATH. "+
				"Please install pandoc https://pandoc.org/installing.html",
			fromFormat,
		)
		cli.HardError(errorMsg)
		return emptyFile, nil
	}

	cmd := exec.Command(pandocPath, "-f", fromFormat, "-t", toFormat)
	cmd.Stdin = nil
	cmd.Stderr = nil

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return emptyFile, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return emptyFile, err
	}

	if err := cmd.Start(); err != nil {
		return emptyFile, err
	}

	_, err = stdin.Write([]byte(content))
	if err != nil {
		return emptyFile, err
	}

	stdin.Close()

	outputBytes := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			outputBytes = append(outputBytes, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		return emptyFile, err
	}

	return outputBytes, nil
}
