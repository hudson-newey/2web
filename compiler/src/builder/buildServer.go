package builder

import (
	"fmt"
	"path/filepath"
	"regexp"
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

// Server side code comes in two kinds of files:
//
//   - HTTP verb routes ("__get.server.ts", "__post.server.ts", ...): each
//     default exports an express request handler that handles one HTTP
//     method, and is mounted on the generated express server at the route
//     that mirrors the directory the file sits in.
//
//     e.g. "src/api/users/__get.server.ts" handles "GET /api/users".
//
//   - Server modules (any other ".server.ts" file): compiled so that their
//     named exports can be invoked over rpc (see the html outputs, compiled
//     functions, and event reducers that call them), but never mounted as
//     http endpoints. Verb routes are the only entry point for accessing a
//     server endpoint over http.

// Server output is emitted next to the client output with a "-server"
// suffix. e.g. "./dist/" -> "./dist-server/".
func serverOutputDir() string {
	outputPath := cli.GetArgs().OutputPath
	return strings.TrimSuffix(strings.TrimSuffix(outputPath, "/"), "\\") + "-server"
}

// ServerRoute describes one HTTP verb route in the route manifest.
type ServerRoute struct {
	// Route is the url path that the handler is mounted at.
	// e.g. "/api/users"
	Route string `json:"route"`

	// Method is the (lowercase) http method the handler responds to.
	// e.g. "get"
	Method string `json:"method"`

	// File is the path of the compiled handler module, relative to the server
	// output directory.
	// e.g. "api/users/__get.server.js"
	File string `json:"file"`
}

// ServerModule describes a compiled server module whose exported functions
// are callable over rpc (see the rpc endpoint in the server runtime). Server
// modules are never mounted as http endpoints.
type ServerModule struct {
	// File is the path of the compiled module, relative to the server output
	// directory.
	// e.g. "api/users.server.js"
	File string `json:"file"`

	// Rpc lists the names of the exported functions that the generated rpc
	// endpoints call.
	Rpc []string `json:"rpc,omitempty"`
}

var (
	serverRoutesMutex sync.Mutex
	serverRoutes      []ServerRoute
	serverModules     []ServerModule
)

// serverVerbs lists the http methods that a "__<method>.server.ts" route
// file can handle.
var serverVerbs = []string{"get", "post", "put", "patch", "delete", "head", "options"}

// verbRouteForServerScript returns the http method and the url path of a
// verb route file (e.g. "__get.server.ts"), and whether the file is a verb
// route at all.
//
// The route mirrors the directory the file sits in:
//
//	(with an input path of "src/")
//	src/api/users/__get.server.ts -> ("get", "/api/users")
//	src/__post.server.ts          -> ("post", "/")
func verbRouteForServerScript(inputPath string, filePath string) (method string, route string, isVerbRoute bool) {
	base := strings.ToLower(filepath.Base(filePath))

	for _, verb := range serverVerbs {
		for _, extension := range []string{".ts", ".js", ".mjs"} {
			if base != "__"+verb+".server"+extension {
				continue
			}

			directory := filepath.Dir(filePath)

			relativePath, err := filepath.Rel(inputPath, directory)
			if err != nil || relativePath == "." {
				return verb, "/", true
			}

			return verb, "/" + filepath.ToSlash(relativePath), true
		}
	}

	return "", "", false
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

	relativeOutput := serverScriptOutputPath(inputPath, filePath)

	serverRoutesMutex.Lock()
	if method, route, isVerbRoute := verbRouteForServerScript(inputPath, filePath); isVerbRoute {
		// A verb route file is an http endpoint: its default export handles
		// the method the file name declares. Named exports of a verb route
		// are not rpc callable (the route is the entry point).
		serverRoutes = append(serverRoutes, ServerRoute{
			Route:  route,
			Method: method,
			File:   relativeOutput,
		})
	} else {
		// A server module is only reachable over rpc.
		serverModules = append(serverModules, ServerModule{
			File: relativeOutput,
			Rpc:  exportedFunctionNames(string(source)),
		})
	}
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
	modules := serverModules
	serverRoutesMutex.Unlock()

	if len(routes) == 0 && len(modules) == 0 {
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
			"    { \"route\": %q, \"method\": %q, \"file\": %q }%s\n",
			route.Route, route.Method, route.File, separator,
		)
	}
	manifestContent += "  ],\n  \"modules\": [\n"
	for i, module := range modules {
		separator := ","
		if i == len(modules)-1 {
			separator = ""
		}

		rpcFunctions := "[]"
		if len(module.Rpc) > 0 {
			quoted := make([]string, len(module.Rpc))
			for index, name := range module.Rpc {
				quoted[index] = fmt.Sprintf("%q", name)
			}

			rpcFunctions = "[ " + strings.Join(quoted, ", ") + " ]"
		}

		manifestContent += fmt.Sprintf(
			"    { \"file\": %q, \"rpc\": %s }%s\n",
			module.File, rpcFunctions, separator,
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

// exportedFunctionNames scans a server script source for the functions that
// are exported for rpc calls.
//
// The scan is intentionally conservative: only named function declarations and
// arrow function constants are picked up. Default exports are route handlers
// (not rpc targets), and re-exports aren't supported.
func exportedFunctionNames(source string) []string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^\s*export\s+(?:async\s+)?function\s+([A-Za-z_$][\w$]*)`),
		regexp.MustCompile(`export\s+const\s+([A-Za-z_$][\w$]*)\s*=\s*(?:async\s*)?(?:\(|function)`),
	}

	seen := map[string]bool{}
	names := []string{}

	for _, pattern := range patterns {
		for _, match := range pattern.FindAllStringSubmatch(source, -1) {
			name := match[1]

			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}

	return names
}
