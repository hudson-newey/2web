package shell

import (
	"log"
	"os"
	"os/exec"
)

// Executes a shell command while providing feedback
func ExecuteCommand(command ...string) error {
	log.Println(command)

	commandLen := len(command)
	if commandLen == 0 {
		panic("command must have a length greater than 0")
	}

	var cmd *exec.Cmd
	if commandLen > 1 {
		cmd = exec.Command(command[0], command[1:]...)
	} else {
		cmd = exec.Command(command[0])
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
