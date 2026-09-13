package templating

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
	optimizer "hudson-newey/2web/src/compiler/6-optimizer"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/debugger"
	"os"
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

	// Main part of the compiler where we recurse the constructed AST
	recurseAstMarkup(&pageModel, parsedAst)
	recurseAst(&pageModel, parsedAst, parsedAst)

	if !cli.GetArgs().IsolatedPages {
		addRouteAssets(&pageModel)
	}

	// The debug file is only generated for development builds (see site.AfterAll),
	// so collecting the reactive graph is skipped for production builds.
	if !cli.GetArgs().IsProd {
		debugger.AddPageDebugInfo(nodes.CollectDebugInfo(filePath, parsedAst))
	}

	args := cli.GetArgs()
	if args.WithFormatting {
		pageModel.Format()
	}

	// We always optimize last so that even the injected content is optimized.
	optimizer.OptimizePage(&pageModel)

	return pageModel
}

func recurseAstMarkup(page *page.Page, parsedAst nodes.AbstractSyntaxTree) {
	// First pass to establish page content
	// We need this so that when we get to ast nodes that replace page content,
	// it has the full page context.
	for _, node := range parsedAst {
		markupContent := node.MarkupContent()
		page.SetContent(page.Html.Content + markupContent)

		recurseAstMarkup(page, node.Children())
	}
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
func recurseAst(page *page.Page, rootAst nodes.AbstractSyntaxTree, current nodes.AbstractSyntaxTree) {
	// Second pass reactive content
	for _, node := range current {
		nodeContent := node.Content(page, rootAst)
		page.SetContent(nodeContent.HtmlContent.Content)
		page.AddStyle(nodeContent.CssContent)
		page.AddScript(nodeContent.JsContent)
		page.AddTwoScript(nodeContent.TwoScriptContent)

		recurseAst(page, rootAst, node.Children())
	}
}
