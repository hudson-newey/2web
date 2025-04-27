package install

import (
	"fmt"
	"os/exec"
)

func installBunPackage(name string) {
	cmd := exec.Command("bun", "add", name)
	stdout, err := cmd.Output()

	if err != nil {
		errorMsg := fmt.Errorf("failed to install package '%s': %s", name, err)
		panic(errorMsg)
	}

	fmt.Println(string(stdout))
}
