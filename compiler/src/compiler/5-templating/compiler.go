package templating

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
	optimizer "hudson-newey/2web/src/compiler/6-optimizer"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/debugger"
	"os"
	"strings"
)

// This package is called concurrently from the build thread pool, so all page
// compilation state must be local to the Compile call. A package level AST
// would be overwritten by every concurrently compiled page, causing pages to
// render content from other pages' ASTs.
func Compile(filePath string, parsedAst nodes.AbstractSyntaxTree) page.Page {
	// Raw text files should be returned without modification since they are "raw"
	// data formats.
	// Adding additional functionality on top of these formats is nonsensical and
	// would only increase the surface for bugs.
	if txt.IsTxtFile(filePath) {
		fileContent, _ := os.ReadFile(filePath)
		pageModel := page.NewPage()
		pageModel.SetContent(string(fileContent))
		return pageModel
	}

	pageModel := page.NewPage()
	pageModel.InputPath = filePath

	// Pre-resolve the reactive relationships of the page once. Dependency
	// queries during compilation are answered from this index instead of
	// walking the whole page AST per variable per candidate.
	reactiveIndex := nodes.BuildReactiveIndex(parsedAst)

	// Pre-allocate the runtime variable names of the reactive runtime. Every
	// variable that needs a runtime representation gets one before the
	// reactivity pass so that computed variables can reference the runtime
	// variables they are computed from regardless of compilation order.
	for _, variable := range reactiveIndex.Variables {
		needsRuntime := reactiveIndex.IsRuntime(variable) ||
			reactiveIndex.HasDerivedDependents(variable) ||
			(reactiveIndex.IsDerived(variable) && reactiveIndex.HasRuntimeDependency(variable))

		if needsRuntime {
			reactiveIndex.RegisterRuntimeVariable(
				variable.Selector(),
				pageModel.Ids.CreateVariableName(),
			)
		}
	}

	// Main part of the compiler where we recurse the constructed AST
	recurseAstMarkup(&pageModel, parsedAst)
	recurseAst(&pageModel, parsedAst, reactiveIndex)

	if !cli.GetArgs().IsolatedPages {
		addRouteAssets(&pageModel)
	}

	// Emit the shared reactive runtime: the compiled wiring of every reactive
	// variable, in dependency order, in one shared scope so that computed
	// variables can reference the runtime variables they are computed from.
	pageModel.EmitReactiveRuntime()

	// The debug file is only generated for development builds (see site.AfterAll),
	// so collecting the reactive graph is skipped for production builds.
	if !cli.GetArgs().IsProd {
		debugger.AddPageDebugInfo(nodes.CollectDebugInfo(filePath, reactiveIndex))
	}

	args := cli.GetArgs()
	if args.WithFormatting {
		pageModel.Format()
	}

	// We always optimize last so that even the injected content is optimized.
	optimizer.OptimizePage(&pageModel)

	return pageModel
}

// recurseAstMarkup runs the first pass to establish page content.
//
// The markup is accumulated in a strings.Builder and written to the page once.
// Appending every node's markup to the page content string would be quadratic
// in the number of nodes (and in the size of the page), because every
// concatenation copies the whole document so far.
//
// This pass must not read the page content: it only concatenates node markup.
func recurseAstMarkup(page *page.Page, parsedAst nodes.AbstractSyntaxTree) {
	var markup strings.Builder

	var walk func(ast nodes.AbstractSyntaxTree)
	walk = func(ast nodes.AbstractSyntaxTree) {
		for _, node := range ast {
			markup.WriteString(node.MarkupContent())
			walk(node.Children())
		}
	}

	walk(parsedAst)
	page.SetContent(markup.String())
}

// recurseAst runs the reactive compilation pass over every node in the AST.
//
// Every node is handed the ROOT AST (rootAst) rather than the subtree it is
// currently in, because reactive nodes need the whole page's AST to find their
// dependencies. For example, a reactive variable declared inside a
// <script compiled> block depends on event and property nodes that are spread
// across the whole page.
//
// The root is threaded through the recursion instead of being stashed in a
// package global because this package is called concurrently from the build
// thread pool (a global would be overwritten by every concurrently compiled
// page, causing pages to render content from other pages' ASTs).
func recurseAst(page *page.Page, rootAst nodes.AbstractSyntaxTree, index *nodes.ReactiveIndex) {
	// Second pass reactive content
	for _, node := range rootAst {
		nodeContent := node.Content(page, index)
		page.SetContent(nodeContent.HtmlContent.Content)
		page.AddStyle(nodeContent.CssContent)
		page.AddScript(nodeContent.JsContent)
		page.AddTwoScript(nodeContent.TwoScriptContent)

		recurseAst(page, node.Children(), index)
	}
}
