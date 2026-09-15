//go:build windows

package processes

import (
	"os/exec"
	"strconv"
)

// processAlive reports whether a process with the given pid exists.
//
// Windows doesn't support the empty signal probe, so the process list is
// consulted instead.
func processAlive(pid int) bool {
	cmd := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(pid), "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(output) > 0 && !containsNoTasks(output)
}

func containsNoTasks(output []byte) bool {
	return len(output) > 0 && string(output[:]) == "" || len(output) == 0
}
