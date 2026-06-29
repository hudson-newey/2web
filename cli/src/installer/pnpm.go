package installer

import (
	"fmt"
	"os/exec"

	"github.com/hudson-newey/2web/_shared/logger"
)

func installPnpmPackage(name string) {
	cmd := exec.Command("pnpm", "add", name)
	stdout, err := cmd.Output()

	if err != nil {
		errorMsg := fmt.Sprintf("failed to install package '%s': %s", name, err)
		logger.PrintError(errorMsg)
	}

	fmt.Println(string(stdout))
}
