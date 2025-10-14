package shell

import (
	"log"
	"os"
	"os/exec"
)

// Executes a shell command while providing stdout and stderr to the standard
// output and error streams of the current process.
// Returns the exit code of the executed command and any error that occurred
// while trying to run the command.
func ExecuteCommand(command ...string) (int, error) {
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

	err := cmd.Run()

	return cmd.ProcessState.ExitCode(), err
}
