package cli

import "log"

func PrintBuildLog(message string) {
	if !*GetArgs().IsSilent {
		log.Println(message)
	}
}
