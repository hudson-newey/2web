package commands

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/builders"
	"github.com/hudson-newey/2web-cli/src/builders/build"
	"github.com/hudson-newey/2web-cli/src/processes"
	"github.com/hudson-newey/2web-cli/src/ssr"
	"github.com/hudson-newey/2web/_shared/logger"
)

// The "server" command group runs and stops the compiled server routes.
//
// Server routes are compiled by the 2web compiler from `.server.ts` files into
// the server output directory (e.g. ./dist-server/), where they are mounted on
// a generated express server (see the route manifest, routes.json).
func serverCommand(programName string, command string, args []string) {
	argsLen := len(args)
	if argsLen < 3 {
		printServerUsage(programName, command)
		return
	}

	subCommand := args[2]

	// The server runtime is executed from the server output directory so that
	// its imports (the compiled handlers and @two-web/kit) resolve.
	serverEntry := ssr.ServerOutputDir + "main.js"

	switch subCommand {
	case "run":
		runServerRuntime(serverEntry, args)
	case "start":
		startServerRuntime(serverEntry, args)
	case "stop":
		stopServerRuntime()
	case "status":
		serverStatus()
	default:
		errorMsg := fmt.Sprintf("unknown sub command: '%s'", subCommand)
		logger.PrintError(errorMsg)
		printServerUsage(programName, command)
	}
}

// Compiles the solution (client + server output) before the server runtime is
// started, so that the server routes are up to date.
func ensureServerRoutesBuilt(args []string) {
	if _, err := os.Stat(ssr.ServerRouteManifest); err == nil {
		return
	}

	build.BuildSolution(args)
}

func serverRuntimeEnvironment(args []string) map[string]string {
	port := serverPort(args)
	if port == "" {
		return map[string]string{}
	}

	return map[string]string{"PORT": port}
}

// serverPort extracts the port from the optional --port flag.
func serverPort(args []string) string {
	for index, arg := range args {
		if arg == "--port" && index+1 < len(args) {
			return args[index+1]
		}
	}

	return ""
}

func runServerRuntime(serverEntry string, args []string) {
	ensureServerRoutesBuilt(args)

	if _, err := os.Stat(serverEntry); err != nil {
		logger.PrintError(
			fmt.Sprintf(
				"no compiled server routes found ('%s' is missing). Add a '.server.ts' file to your project and rebuild.",
				serverEntry,
			),
		)
		return
	}

	if err := processes.Run(serverEntry, serverRuntimeEnvironment(args)); err != nil {
		logger.PrintError(fmt.Sprintf("server exited with an error: %v", err))
	}
}

func startServerRuntime(serverEntry string, args []string) {
	ensureServerRoutesBuilt(args)

	if _, err := os.Stat(serverEntry); err != nil {
		logger.PrintError(
			fmt.Sprintf(
				"no compiled server routes found ('%s' is missing). Add a '.server.ts' file to your project and rebuild.",
				serverEntry,
			),
		)
		return
	}

	pid, err := processes.Start(serverEntry, serverRuntimeEnvironment(args))
	if err != nil {
		logger.PrintError(err.Error())
		return
	}

	logger.Println(fmt.Sprintf("Started the 2web server (pid %d)", pid))
	logger.Println("Stop it with '2web server stop'")
}

func stopServerRuntime() {
	if err := processes.Stop(); err != nil {
		logger.PrintError(err.Error())
	}
}

func serverStatus() {
	running, pid, _ := processes.Status()

	if running {
		logger.Println(fmt.Sprintf("The 2web server is running (pid %d)", pid))
		return
	}

	logger.Println("The 2web server is not running")
}

func printServerUsage(programName string, command string) {
	fmt.Printf(`invalid arguments:
	expected: %s %s <sub_command>

  %s %s run    Compiles and runs the server routes in the foreground
  %s %s start  Compiles and runs the server routes in the background
  %s %s stop   Stops the background server
  %s %s status Reports whether the server is running
`, programName, command, programName, command, programName, command, programName, command, programName, command)
}

// hasServerRoutesSource returns whether the solution contains any server
// script sources (which the compiler turns into server routes).
func hasServerRoutesSource() bool {
	for _, entry := range entryCandidates() {
		if directoryContainsServerScript(entry) {
			return true
		}
	}

	return false
}

func entryCandidates() []string {
	candidates := []string{"./src/", "./server/"}
	return candidates
}

// _ kept for documentation: the compiler routes .server.ts files into the
// server output directory (see the 2web compiler's buildServer.go).
var _ = builders.EntryTargets

// directoryContainsServerScript walks the directory tree looking for a
// `.server.ts` (or `.server.js`) file.
func directoryContainsServerScript(directory string) bool {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		fullPath := directory + "/" + entry.Name()

		if entry.IsDir() {
			if directoryContainsServerScript(fullPath) {
				return true
			}

			continue
		}

		if len(fullPath) > 10 && fullPath[len(fullPath)-10:] == ".server.ts" {
			return true
		}

		if len(fullPath) > 10 && fullPath[len(fullPath)-10:] == ".server.js" {
			return true
		}
	}

	return false
}

var _ = build.BuildSolution
var _ = stopServerRuntime
var _ = serverStatus
var _ = serverRuntimeEnvironment
var _ = hasServerRoutesSource
var _ = directoryContainsServerScript
