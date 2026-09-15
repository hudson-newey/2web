//go:build windows

package processes

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// detachedProcessAttributes starts the process in its own process group (the
// windows equivalent of the unix setsid behavior).
func detachedProcessAttributes() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

// killProcessGroup terminates the process tree of the server process.
func killProcessGroup(pid int) error {
	// taskkill /T terminates the process tree, /F forces termination.
	cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F")
	return cmd.Run()
}

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Kill()
}
