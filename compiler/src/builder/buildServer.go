package builder

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/filesystem"
	"hudson-newey/2web/src/models"
)

// A server route is a `.server.ts` (or `.server.js`) file. Each file is
// compiled into a standalone node module that default exports an express
// request handler, and is mounted on the generated express server at the
// route that mirrors its position in the source tree.
//
// e.g. "src/api/users.server.ts" is served at "/api/users".

// Server output is emitted next to the client output with a "-server"
// suffix. e.g. "./dist/" -> "./dist-server/".
func serverOutputDir() string {
	outputPath := cli.GetArgs().OutputPath
	return strings.TrimSuffix(strings.TrimSuffix(outputPath, "/"), "\\") + "-server"
}

// ServerRoute describes a compiled server route in the route manifest.
type ServerRoute struct {
	// Route is the url path that the handler is mounted at.
	// e.g. "/api/users"
	Route string `json:"route"`

	// File is the path of the compiled handler module, relative to the server
	// output directory.
	// e.g. "api/users.server.js"
	File string `json:"file"`
}

var (
	serverRoutesMutex sync.Mutex
	serverRoutes      []ServerRoute
)

// routeForServerScript derives the url path of a server script from its
// position in the source tree.
//
// e.g. (with an input path of "src/"):
//
//	src/api/users.server.ts -> /api/users
//	src/index.server.ts     -> /
//	src/api/index.server.ts -> /api
func routeForServerScript(inputPath string, filePath string) string {
	relativePath, err := filepath.Rel(inputPath, filePath)
	if err != nil {
		// Fall back to the file name when the path isn't inside the input
		// path (e.g. a file passed directly to the compiler).
		relativePath = filepath.Base(filePath)
	}

	route := "/" + filepath.ToSlash(relativePath)

	// Strip the server script extension. e.g. "api/users.server.ts" -> "api/users"
	route = strings.TrimSuffix(route, ".server.ts")
	route = strings.TrimSuffix(route, ".server.js")
	route = strings.TrimSuffix(route, ".server.mjs")

	// Directory index files are served at the directory route.
	route = strings.TrimSuffix(route, "index")
	if route == "" {
		route = "/"
	}

	return route
}

// serverScriptOutputPath returns the compiled handler output path for a
// server script, relative to the server output directory.
//
// e.g. "src/api/users.server.ts" -> "api/users.server.js"
func serverScriptOutputPath(inputPath string, filePath string) string {
	relativePath, err := filepath.Rel(inputPath, filePath)
	if err != nil {
		relativePath = filepath.Base(filePath)
	}

	return filepath.ToSlash(strings.TrimSuffix(relativePath, filepath.Ext(relativePath)) + ".js")
}

// buildServerRoute compiles a `.server.ts` (or `.server.js`) file into a
// standalone node module in the server output directory.
//
// Server code is bundled for the node runtime: dependencies from the project's
// node_modules stay external (so that server only dependencies, database
// drivers, etc... aren't duplicated into the bundle), while the user's own
// modules are bundled in.
func buildServerRoute(inputPath string, filePath string) {
	source, err := filesystem.ReadFile(filePath)
	if err != nil {
		readError := models.NewError(
			"failed to read server script: "+err.Error(),
			filePath,
			lexer.StartingPosition,
		)

		documentErrors.AddErrors(&readError)
		cli.PrintBuildLog("\t- " + filePath + " \033[31m(ERROR)\033[0m")
		return
	}

	workingDir, _ := filepath.Abs(inputPath)

	// Imports are resolved relative to the directory of the server script
	// (not the compiler's working directory), so that a script can import its
	// siblings with relative paths.
	resolveDir, _ := filepath.Abs(filepath.Dir(filePath))

	esbuildOutput := api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   string(source),
			Loader:     api.LoaderTS,
			ResolveDir: resolveDir,
			Sourcefile: filepath.Base(filePath),
		},
		Platform:         api.PlatformNode,
		Format:           api.FormatESModule,
		AbsWorkingDir:    workingDir,
		Sourcemap:        api.SourceMapNone,
		Bundle:           true,
		Packages:         api.PackagesExternal,
		MinifyWhitespace: cli.GetArgs().IsProd,
	})

	bundledContent := ""
	for _, outputFile := range esbuildOutput.OutputFiles {
		bundledContent += string(outputFile.Contents)
	}

	for _, buildError := range esbuildOutput.Errors {
		position := lexer.StartingPosition
		if buildError.Location != nil {
			position = lexer.Position{
				Row: int(buildError.Location.Line),
				Col: int(buildError.Location.Column),
			}
		}

		errorModel := models.NewError(buildError.Text, filePath, position)
		documentErrors.AddErrors(&errorModel)
	}

	if len(esbuildOutput.Errors) > 0 || bundledContent == "" {
		cli.PrintBuildLog("\t- " + filePath + " \033[31m(ERROR)\033[0m")
		return
	}

	outputPath := filepath.Join(serverOutputDir(), serverScriptOutputPath(inputPath, filePath))
	filesystem.WriteFile([]byte(bundledContent), outputPath)

	route := routeForServerScript(inputPath, filePath)

	serverRoutesMutex.Lock()
	serverRoutes = append(serverRoutes, ServerRoute{
		Route: route,
		File:  serverScriptOutputPath(inputPath, filePath),
	})
	serverRoutesMutex.Unlock()

	cli.PrintBuildLog("\t- " + filePath + " \033[35m(server)\033[0m")
}

// FlushServerRoutes writes the compiled server route manifest and the server
// entry point into the server output directory.
//
// The manifest is consumed by the server runtime (see @two-web/kit/ssr) to
// mount the compiled handlers on the generated express server.
//
// It must be called after all server routes have been compiled.
func FlushServerRoutes() {
	serverRoutesMutex.Lock()
	routes := serverRoutes
	serverRoutesMutex.Unlock()

	if len(routes) == 0 {
		return
	}

	outputDirectory := serverOutputDir()

	manifestContent := "{\n  \"routes\": [\n"
	for i, route := range routes {
		separator := ","
		if i == len(routes)-1 {
			separator = ""
		}

		manifestContent += fmt.Sprintf(
			"    { \"route\": %q, \"file\": %q }%s\n",
			route.Route, route.File, separator,
		)
	}
	manifestContent += "  ]\n}\n"

	if err := filesystem.WriteFileSync([]byte(manifestContent), filepath.Join(outputDirectory, "routes.json")); err != nil {
		manifestError := models.NewError(
			"failed to write the server route manifest: "+err.Error(),
			outputDirectory,
			lexer.StartingPosition,
		)

		documentErrors.AddErrors(&manifestError)
		return
	}

	// The server entry point boots the generated express server. It imports
	// the server runtime from @two-web/kit, which is a dependency of every
	// 2web project that compiles server routes (the same dependency that the
	// "2web template ssr" template uses).
	entryContent := `import path from "node:path";
import { fileURLToPath } from "node:url";
import { startRouteServer } from "@two-web/kit/ssr";

startRouteServer({
  serverDir: path.dirname(fileURLToPath(import.meta.url)),
});
`

	if err := filesystem.WriteFileSync([]byte(entryContent), filepath.Join(outputDirectory, "main.js")); err != nil {
		entryError := models.NewError(
			"failed to write the server entry point: "+err.Error(),
			outputDirectory,
			lexer.StartingPosition,
		)

		documentErrors.AddErrors(&entryError)
	}
}

// IsServerScript returns whether the given file is a server script.
func IsServerScript(filePath string) bool {
	return content.IsServerTarget(filePath) &&
		(strings.HasSuffix(filePath, ".ts") ||
			strings.HasSuffix(filePath, ".js") ||
			strings.HasSuffix(filePath, ".mjs"))
}
