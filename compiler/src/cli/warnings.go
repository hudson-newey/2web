package cli

import (
	"github.com/hudson-newey/2web/_shared/logger"
)

func PrintWarning(message string) {
	if GetArgs().IsSilent {
		return
	}

	logger.PrintWarning(message)
}
