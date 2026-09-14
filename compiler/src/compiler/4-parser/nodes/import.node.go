package nodes

import (
	"errors"
	"fmt"
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
	"hudson-newey/2web/src/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// import importName from "importPath";
func NewscriptImportNode(lexNodes []*lexer.V2LexNode, context *ParseContext) Node {
	importNameNode, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 1)
	if err != nil {
		return context.DegradedNode(
			"import statement is missing an imported name. Components are imported with 'import Name from \"path\";'",
			lexNodes,
		)
	}

	importNameNode.Content = strings.TrimSpace(importNameNode.Content)

	importPathNode, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 2)
	if err != nil {
		return context.DegradedNode(
			"import statement is missing an import path. Components are imported with 'import Name from \"path\";'",
			lexNodes,
		)
	}

	quoteRe := regexp.MustCompile(`"(.*?)"|'(.*?)'`)
	quotedPath := quoteRe.FindAllStringSubmatch(importPathNode.Content, -1)
	if len(quotedPath) == 0 {
		return context.DegradedNode(
			"import path must be a quoted string (e.g. 'import Name from \"components/header.component.html\";')",
			lexNodes,
		)
	}

	// The path is captured in the first or second sub match depending on the
	// quote style used.
	importPath := quotedPath[0][1]
	if importPath == "" {
		importPath = quotedPath[0][2]
	}

	importPathNode.Content = importPath

	return &scriptImportNode{
		importName: importNameNode.Content,
		importPath: strings.TrimSpace(importPathNode.Content),
	}
}

type scriptImportNode struct {
	// What was the alias that the import was given?
	importName string

	// What is the path of the import?
	importPath string

	children AbstractSyntaxTree
}

func (m *scriptImportNode) Type() string {
	return "importNode"
}

func (m *scriptImportNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *scriptImportNode) MarkupContent() string {
	return ""
}

func (m *scriptImportNode) Content(page *page.Page, index *ReactiveIndex) NodeContent {
	// Check that the imported path really exists.
	hostDirectory := filepath.Dir(page.InputPath)
	componentPath := filepath.Join(hostDirectory, m.importPath)

	// Compile time virtual functions (the interop types: $props(), $uid(),
	// $env(), $readFile()) are evaluated by the preprocessor, which also
	// removes their import statement. If one still reaches this node, it
	// contributes nothing to the page rather than erroring.
	if strings.Contains(m.importPath, "interop.types") {
		return NodeContent{
			HtmlContent:      page.Html,
			JsContent:        javascript.NewJsFile(),
			CssContent:       css.NewCssFile(),
			TwoScriptContent: twoscript.NewTwoScriptFile(),
		}
	}

	// Imports of server scripts (.server.ts files) generate async client
	// functions that pass through to generated rpc endpoints. The imported
	// functions can then be called directly from reactive event listeners.
	if isServerScriptPath(componentPath) {
		if !strings.HasPrefix(strings.TrimSpace(m.importName), "{") {
			// A bare import names the route handler's default export (a bare
			// import is also how html components are imported). The default
			// export isn't rpc callable, so only named imports make sense.
			defaultImportError := models.NewError(
				fmt.Sprintf(
					"server scripts are imported for rpc calls with named imports (e.g. 'import { myFunction } from \"%s\";')",
					m.importPath,
				),
				page.InputPath,
				lexer.Position{},
			)

			page.Errors.AddError(&defaultImportError)

			return NodeContent{
				HtmlContent:      page.Html,
				JsContent:        javascript.NewJsFile(),
				CssContent:       css.NewCssFile(),
				TwoScriptContent: twoscript.NewTwoScriptFile(),
			}
		}

		m.generateRpcPassthroughs(page, index, componentPath)
		return NodeContent{
			HtmlContent:      page.Html,
			JsContent:        javascript.NewJsFile(),
			CssContent:       css.NewCssFile(),
			TwoScriptContent: twoscript.NewTwoScriptFile(),
		}
	}

	componentContent, err := os.ReadFile(componentPath)
	if errors.Is(err, os.ErrNotExist) {
		msg := fmt.Sprintf("error importing file. Could not find file: '%s'", m.importPath)
		pageErr := models.NewError(msg, page.InputPath, lexer.Position{})
		page.Errors.AddError(&pageErr)
	}

	// Replace all of the selectors with the components content.
	selector := fmt.Sprintf("<%s />", m.importName)
	HtmlContent := strings.ReplaceAll(page.Html.Content, selector, string(componentContent))

	return NodeContent{
		HtmlContent:      html.FromContent(HtmlContent),
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

func (m *scriptImportNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *scriptImportNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

// isServerScriptPath returns whether the path points at a server script.
func isServerScriptPath(filePath string) bool {
	return strings.HasSuffix(filePath, ".server.ts") ||
		strings.HasSuffix(filePath, ".server.js") ||
		strings.HasSuffix(filePath, ".server.mjs")
}

// generateRpcPassthroughs generates the client side async functions that pass
// through to the rpc endpoints of the imported server functions.
//
// The generated functions are appended to the page's shared reactive runtime
// (in the same scope as the reactive event handlers), so that an event
// listener calling the imported function directly calls the passthrough.
//
// e.g. 'import { greet } from "./api/users.server.ts";' generates:
//
//	async function greet(...greetArgs) {
//		const response = await fetch("/_2web/rpc/api/users.server.js/greet", {
//			method: "POST",
//			headers: { "content-type": "application/json" },
//			body: JSON.stringify(greetArgs),
//		});
//
//		return response.text();
//	}
func (m *scriptImportNode) generateRpcPassthroughs(
	pageModel *page.Page,
	index *ReactiveIndex,
	serverPath string,
) {
	for _, importName := range m.rpcFunctionNames() {
		if index.HasRpcFunction(importName) {
			// The same function imported by two script blocks generates a
			// single passthrough (duplicate declarations in the shared runtime
			// scope would be a syntax error).
			continue
		}

		index.RegisterRpcFunction(importName)

		module := rpcModulePath(serverPath)
		endpoint := "/_2web/rpc/" + module + "/" + importName
		arguments := importName + "Args"

		passthrough := fmt.Sprintf(
			`async function %s(...%s) {
	const __2_rpc_response = await fetch(%q, {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(%s),
	});

	return __2_rpc_response.text();
}`,
			importName, arguments, endpoint, arguments,
		)

		pageModel.AppendReactiveRuntime(page.ReactiveRuntimeChunk{
			Variable: "rpc:" + importName,
			Content:  passthrough,
		})
	}
}

// rpcFunctionNames parses the imported names of an import statement.
//
// Both import forms are supported:
//
//	import { greet, other } from "..." -> ["greet", "other"]
//	import greet from "..."            -> ["greet"]
func (m *scriptImportNode) rpcFunctionNames() []string {
	trimmed := strings.TrimSpace(m.importName)

	if !strings.HasPrefix(trimmed, "{") {
		return []string{trimmed}
	}

	trimmed = strings.TrimSuffix(strings.TrimPrefix(trimmed, "{"), "}")

	names := []string{}
	for _, name := range strings.Split(trimmed, ",") {
		name = strings.TrimSpace(name)

		if name != "" {
			names = append(names, name)
		}
	}

	return names
}

// rpcModulePath returns the path of the compiled handler module for a server
// script, relative to the server output directory.
//
// This must mirror the output layout of the server build (see the server
// script output path in the builder): both paths are derived from the
// compiler's input path, so that the rpc url addresses the compiled module.
//
// e.g. (with an input path of "src/") "src/api/users.server.ts" is compiled to
// "<server output>/api/users.server.js" and addressed at
// "/_2web/rpc/api/users.server.js".
func rpcModulePath(serverPath string) string {
	inputRoot := cli.GetArgs().InputPath

	relativePath, err := filepath.Rel(inputRoot, serverPath)
	if err != nil {
		// Fall back to the file name when the path isn't inside the input
		// path (the same fallback the server build uses).
		relativePath = filepath.Base(serverPath)
	}

	extension := filepath.Ext(relativePath)

	return filepath.ToSlash(strings.TrimSuffix(relativePath, extension) + ".js")
}
