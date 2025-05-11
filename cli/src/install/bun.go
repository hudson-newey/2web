package install

import (
	"fmt"
	"os/exec"

	"github.com/hudson-newey/2web-cli/src/cli"
)

func installBunPackage(name string) {
	cmd := exec.Command("bun", "add", name)
	stdout, err := cmd.Output()

	if err != nil {
		errorMsg := fmt.Sprintf("failed to install package '%s': %s", name, err)
		cli.PrintError(2, errorMsg)
	}

	fmt.Println(string(stdout))
}
