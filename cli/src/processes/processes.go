package processes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/hudson-newey/2web/_shared/logger"
)

// Process management for the compiled server runtime.
//
// The process is started detached (in its own process group) so that it
// outlives the cli invocation that started it, and is stopped by killing its
// whole process group.

const (
	// The directory that process metadata (pid and log files) is written into.
	stateDir = ".2web"

	PidFile = stateDir + "/server.pid"
	LogFile = stateDir + "/server.log"
)

// Start launches the compiled server runtime in the background.
//
// The entry point is the compiled server entry (e.g. dist-server/main.js)
// and env is appended to the child process environment (e.g. PORT=3000).
func Start(entryPoint string, env map[string]string) (int, error) {
	if IsRunning() {
		return 0, errors.New("the 2web server is already running (stop it with '2web server stop')")
	}

	runtime, args, err := runtimeInvocation(entryPoint)
	if err != nil {
		return 0, err
	}

	cmd := exec.Command(runtime, args...)
	cmd.Env = append(os.Environ(), environmentList(env)...)

	if err := os.MkdirAll(path.Dir(PidFile), os.ModePerm); err != nil {
		return 0, fmt.Errorf("failed to create the state directory: %w", err)
	}

	logFile, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open the server log file: %w", err)
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = detachedProcessAttributes()

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start the server runtime: %w", err)
	}

	pid := cmd.Process.Pid

	if err := os.WriteFile(PidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return pid, fmt.Errorf("failed to write the server pid file: %w", err)
	}

	// Detached processes are reaped by init; don't wait for the child here.
	go func() {
		_ = cmd.Wait()
	}()

	return pid, nil
}

// Run launches the compiled server runtime in the foreground and waits for it
// to exit.
func Run(entryPoint string, env map[string]string) error {
	runtime, args, err := runtimeInvocation(entryPoint)
	if err != nil {
		return err
	}

	cmd := exec.Command(runtime, args...)
	cmd.Env = append(os.Environ(), environmentList(env)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Stop terminates a running detached server process.
func Stop() error {
	pid, err := readPidFile()
	if err != nil {
		return err
	}

	if !isProcessRunning(pid) {
		_ = os.Remove(PidFile)
		return errors.New("the 2web server is not running (stale pid file removed)")
	}

	// Kill the whole process group so that child processes of the server don't
	// outlive the stop command. Fall back to a single process kill when group
	// termination isn't supported on this platform.
	if err := killProcessGroup(pid); err != nil {
		return fmt.Errorf("failed to stop the server (pid %d): %w", pid, err)
	}

	_ = os.Remove(PidFile)

	logger.Println(fmt.Sprintf("Stopped the 2web server (pid %d)", pid))
	return nil
}

// Status reports whether the compiled server is running.
func Status() (bool, int, error) {
	pid, err := readPidFile()
	if err != nil {
		return false, 0, nil
	}

	return isProcessRunning(pid), pid, nil
}

// IsRunning reports whether a detached server process is running.
func IsRunning() bool {
	running, _, _ := Status()
	return running
}

func readPidFile() (int, error) {
	content, err := os.ReadFile(PidFile)
	if err != nil {
		return 0, fmt.Errorf("the 2web server is not running")
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return 0, fmt.Errorf("the server pid file is corrupt")
	}

	return pid, nil
}

func isProcessRunning(pid int) bool {
	return processAlive(pid)
}

func runtimeInvocation(entryPoint string) (string, []string, error) {
	// The compiled server entry point is an es module; node runs it directly
	// and bun runs it with full permissions.
	runtime, err := exec.LookPath("node")
	if err == nil {
		return runtime, []string{entryPoint}, nil
	}

	bun, bunErr := exec.LookPath("bun")
	if bunErr == nil {
		return bun, []string{entryPoint}, nil
	}

	return "", nil, errors.New("could not find node (or bun) in the PATH to run the server runtime")
}

func environmentList(env map[string]string) []string {
	list := []string{}
	for key, value := range env {
		list = append(list, key+"="+value)
	}

	return list
}
