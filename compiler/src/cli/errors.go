package cli

import (
	"fmt"
	"hudson-newey/2web/src/models"
	"strings"
)

func PrintError(errorModel models.Error) {
	if *GetArgs().IsSilent {
		return
	}

	// Because all error messages are indented by one tab width, I replace all of
	// I add a tab character after every new line.
	formattedErrorMessage := strings.ReplaceAll(errorModel.Message, "\n", "\n\t")

	formattedLineNumber := ""
	if errorModel.Position.LineNumber != 0 {
		formattedLineNumber = fmt.Sprintf(
			" (ln: %d, col: %d)",
			errorModel.Position.LineNumber,
			errorModel.Position.CharNumber,
		)
	}

	// Add a double new line at the start so that the error message is
	// emphasized in the compiler logs.
	fmt.Printf("\n\033[31m[Error] \033[36m%s\033[33m%s\033[0m\n", errorModel.FilePath, formattedLineNumber)
	fmt.Printf("\t%s\n", formattedErrorMessage)
}
