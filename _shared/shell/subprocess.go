package shell

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Executes a shell command while providing stdout and stderr to the standard
// output and error streams of the current process.
//
// This function will continue the execution of the current program, allowing
// the caller to continue execution.
// The sub-process will be strongly linked to the caller.
// Once the caller program exits, the sub-process will also be killed.
// Killing or crashing the sub-process will also cause the caller program to
// exit to prevent undefined behavior.
func ExecuteSubProcess(command ...string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	commandLen := len(command)
	if commandLen == 0 {
		panic("command must have a length greater than 0")
	}

	var cmd *exec.Cmd
	if commandLen > 1 {
		cmd = exec.CommandContext(ctx, command[0], command...)
	} else {
		cmd = exec.CommandContext(ctx, command[0])
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Ensure the child process dies if the parent dies
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}

	cmd.Start()
}
