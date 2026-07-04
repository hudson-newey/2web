package cli

import "github.com/hudson-newey/2web/_shared/logger"

func PrintBuildLog(message string) {
	if !GetArgs().IsSilent {
		logger.Println(message)
	}
}
