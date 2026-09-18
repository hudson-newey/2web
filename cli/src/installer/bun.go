package installer

import (
	"fmt"

	"github.com/hudson-newey/2web/_shared/logger"
	"github.com/hudson-newey/2web/_shared/shell"
)

func installBunPackage(name string) {
	_, err := shell.ExecuteCommand("bun", "add", name)

	if err != nil {
		errorMsg := fmt.Sprintf("failed to install package '%s': %s", name, err)
		logger.PrintError(errorMsg)
	}
}
