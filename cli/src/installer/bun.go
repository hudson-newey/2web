package installer

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/shell"
	"github.com/hudson-newey/2web/_shared/logger"
)

func installBunPackage(name string) {
	_, err := shell.ExecuteCommand("bun", "add", name)

	if err != nil {
		errorMsg := fmt.Sprintf("failed to install package '%s': %s", name, err)
		logger.PrintError(errorMsg)
	}
}
