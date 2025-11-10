package builder

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/models"
	"io"
	"os"
)

func getContent(inputPath string) *[]byte {
	rawData, err := getInputContent(inputPath)
	if err != nil {
		inputError := models.NewError(
			fmt.Sprintf("Failed to read file: %s\n%s", inputPath, err.Error()),
			inputPath,
			lexer.StartingPosition,
		)

		documentErrors.AddErrors(&inputError)

		// If there was an error reading the file, we return an empty byte slice
		// instead of returning a potentially partially read/corrupted byte slice
		// from the failed read attempt.
		emptyFile := []byte{}
		return &emptyFile
	}

	return &rawData
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
