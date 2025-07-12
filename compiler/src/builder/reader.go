package builder

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
	"io"
	"os"
)

func getContent(inputPath string) string {
	rawData, err := getInputContent(inputPath)
	if err != nil {
		rawData = []byte{}
		documentErrors.AddErrors(models.Error{
			FilePath: inputPath,
			Message:  fmt.Sprintf("Failed to read file: %s\n%s", inputPath, err.Error()),
			Position: models.Position{
				Line: 1,
			},
		})
	}

	return string(rawData)
}

func getInputContent(inputPath string) ([]byte, error) {
	if !*cli.GetArgs().FromStdin {
		return os.ReadFile(inputPath)
	}

	if !*cli.GetArgs().IsSilent {
		fmt.Println("Prompting STDIN for file:", inputPath)
	}

	return io.ReadAll(os.Stdin)
}
