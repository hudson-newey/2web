package devserver

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"

	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/build"
	"github.com/hudson-newey/2web-cli/src/builders/configs"
	"github.com/hudson-newey/2web-cli/src/packages"
	"github.com/hudson-newey/2web-cli/src/runner"
	"github.com/hudson-newey/2web-cli/src/server"
	"github.com/hudson-newey/2web-cli/src/ssr"
	"github.com/hudson-newey/2web/_shared/logger"
)

func ServeSolution(args []string) {
	// Server routes (.server.ts files) are compiled into the server output
	// directory before the dev server starts, so that they can be mounted.
	// The ssr server also serves the compiled client output, so the solution
	// is built whenever it exists too.
	hasRoutes := hasServerSources()
	hasSsr := ssr.HasSsrTarget()
	if hasRoutes || hasSsr {
		build.BuildSolution(args)
	}

	// The ssr server mounts the compiled server routes itself (see the ssr
	// template's runServer), so starting the standalone server runtime would
	// double mount every route on a second port.
	var serverProcess *exec.Cmd
	if hasRoutes && !hasSsr {
		serverProcess = startRouteServerProcess(args)
	}

	switch {
	case hasSsr:
		serveSsr(args)
	case configs.HasViteConfig():
		serveVite(args)
	default:
		serveInbuilt(args)
	}

	// serveSsr and serveVite block until the dev server exits; the server
	// runtime is stopped when it does.
	stopRouteServerProcess(serverProcess)
}

func serveVite(args []string) {
	viteConfig, err := configs.ViteConfigLocation()
	pathTarget := builders.EntryTargets(args)[0]

	// Check that the path target actually exists.
	// If it does not, we want to log a warning.
	if _, err := os.Stat(pathTarget); os.IsNotExist(err) {
		warningMsg := fmt.Sprintf("the specified path target does not exist: '%s'", pathTarget)
		logger.PrintWarning(warningMsg)
	}

	if err == nil {
		packages.ExecutePackage("vite", pathTarget, "--config", viteConfig)
	} else {
		// If there is no vite config, then we want to execute vite without any
		// --config arguments, meaning that Vite should use the default config.
		packages.ExecutePackage("vite", pathTarget)
	}
}

func serveInbuilt(args []string) {
	inPath := builders.EntryTargets(args)[0]
	outPath := builders.OutputTarget(args)

	server.Run(inPath, outPath, server.Options{
		// TODO: Refactor to use go flags package
		WatchFiles: !slices.Contains(args, "--no-watch"),
		AutoReload: !slices.Contains(args, "--no-auto-reload"),
	})
}

func serveSsr(args []string) {
	// The ssr server reads its port from the PORT environment variable (see
	// @two-web/kit/ssr's runServer), so the --port flag is forwarded there.
	for key, value := range serverRuntimeEnvironment(args) {
		os.Setenv(key, value)
	}

	runner.ExecuteScript("./server/ssr.ts")
}

// hasServerSources returns whether the solution contains any .server.ts (or
// .server.js) sources.
func hasServerSources() bool {
	return directoryContainsServerScript("./src/") || directoryContainsServerScript("./server/")
}

// directoryContainsServerScript walks the directory tree looking for a
// `.server.ts` (or `.server.js`) file.
func directoryContainsServerScript(directory string) bool {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		fullPath := path.Join(directory, entry.Name())

		if entry.IsDir() {
			if directoryContainsServerScript(fullPath) {
				return true
			}

			continue
		}

		if isServerScriptPath(fullPath) {
			return true
		}
	}

	return false
}

func isServerScriptPath(filePath string) bool {
	return strings.HasSuffix(filePath, ".server.ts") || strings.HasSuffix(filePath, ".server.js")
}

// startRouteServerProcess starts the compiled server runtime as a child of
// the dev server process.
//
// The child is NOT detached: it shares this process' process group, so a
// single ctrl-c stops both the client dev server and the server runtime.
func startRouteServerProcess(args []string) *exec.Cmd {
	entryPoint := ssr.ServerOutputDir + "main.js"

	if _, err := os.Stat(entryPoint); err != nil {
		logger.PrintWarning(
			fmt.Sprintf("server routes were compiled but the server entry point is missing ('%s')", entryPoint),
		)
		return nil
	}

	runtime, err := exec.LookPath("node")
	if err != nil {
		runtime, err = exec.LookPath("bun")
		if err != nil {
			logger.PrintWarning("could not find node (or bun) in the PATH to run the server runtime")
			return nil
		}
	}

	cmd := exec.Command(runtime, entryPoint)
	cmd.Env = append(os.Environ(), environmentList(serverRuntimeEnvironment(args))...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.PrintWarning(fmt.Sprintf("failed to start the server runtime: %v", err))
		return nil
	}

	return cmd
}

func stopRouteServerProcess(serverProcess *exec.Cmd) {
	if serverProcess == nil || serverProcess.Process == nil {
		return
	}

	if err := serverProcess.Process.Kill(); err != nil {
		// The process already exited.
		return
	}

	_, _ = serverProcess.Process.Wait()
}

// serverRuntimeEnvironment resolves the environment that the compiled server
// runtime is started with.
func serverRuntimeEnvironment(args []string) map[string]string {
	for index, arg := range args {
		if arg == "--port" && index+1 < len(args) {
			return map[string]string{"PORT": args[index+1]}
		}
	}

	return map[string]string{}
}

// environmentList converts an environment map into the list form that
// exec.Command expects.
func environmentList(env map[string]string) []string {
	list := []string{}
	for key, value := range env {
		list = append(list, key+"="+value)
	}

	return list
}
