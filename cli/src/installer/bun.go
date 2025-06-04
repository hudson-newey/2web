package installer

import (
	"fmt"

	"github.com/hudson-newey/2web-cli/src/cli"
	"github.com/hudson-newey/2web-cli/src/shell"
)

func installBunPackage(name string) {
	err := shell.ExecuteCommand("bun", "add", name)

	if err != nil {
		errorMsg := fmt.Sprintf("failed to install package '%s': %s", name, err)
		cli.PrintError(errorMsg)
	}
}
