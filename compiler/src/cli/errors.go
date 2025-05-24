package cli

import (
	"fmt"
	"hudson-newey/2web/src/models"
)

func PrintError(errorModel models.Error) {
	if *GetArgs().IsSilent {
		return
	}

	// Add a double new line so that the error message is emphasized in the
	// compiler logs.
	fmt.Printf("\n\033[31m[Error]\033[0m \033[36m%s\033[0m\n", errorModel.FilePath)
	fmt.Printf("\t%s\n", errorModel.Message)
}
