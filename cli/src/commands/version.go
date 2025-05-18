package commands

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/constants"
)

func printVersion() {
	fmt.Printf("2web cli version: '%s'\n", constants.TwoWebVersion)
}
