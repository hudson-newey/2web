//go:build !windows

package processes

import (
	"fmt"
	"syscall"
)

// detachedProcessAttributes starts the process in its own process group so
// that stopping the server terminates its whole tree.
func detachedProcessAttributes() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup terminates the process group of the server process.
func killProcessGroup(pid int) error {
	// The process group id of a detached process is the process id of its
	// leader.
	return syscall.Kill(-pid, syscall.SIGTERM)
}

func killProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}

var _ = fmt.Sprintf
