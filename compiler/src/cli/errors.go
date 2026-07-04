package cli

import (
	"fmt"

	"github.com/hudson-newey/2web/_shared/logger"
)

// HardError prints a fatal error message and panics to stop execution.
// If the program is erroring because of a document/source error, you should use
// the document.addError method instead, because it can be used to inject errors
// into the HMR overlay, and emit multiple errors for a single run.
//
// Note that because this function panics, it will exit --watch mode.
// You should only enter this state if a non-recoverable error has occurred,
// E.g. a missing dependency, or a major internal failure.
func HardError(message string) {
	// Add a double new line so that the error messages are emphasized in the
	// compiler logs.
	fmt.Println()
	logger.PrintError(message)
	panic("FATAL ERROR")
}
