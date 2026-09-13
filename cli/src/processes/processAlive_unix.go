//go:build !windows

package processes

import (
	"os"
	"syscall"
)

// processAlive reports whether a process with the given pid exists and can be
// signaled.
//
// Signaling with an empty signal doesn't deliver anything; it only checks that
// the process exists.
func processAlive(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return process.Signal(syscall.Signal(0)) == nil
}
