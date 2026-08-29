package templating

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
	optimizer "hudson-newey/2web/src/compiler/6-optimizer"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/txt"
	"os"
)

var rootAstNode []nodes.Node = []nodes.Node{}

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
	rootAstNode = parsedAst
	recurseAstMarkup(&pageModel, parsedAst)
	recurseAst(&pageModel, parsedAst)

	if !cli.GetArgs().IsolatedPages {
		addRouteAssets(&pageModel)
	}

	args := cli.GetArgs()
	if args.WithFormatting {
		pageModel.Format()
	}

	// We always optimize last so that even the injected content is optimized.
	if args.IsProd {
		if args.WithFormatting {
			cli.PrintWarning("Ignoring '--format' because '--production' was specified")
		}

		optimizer.OptimizePage(&pageModel)
	}

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

func recurseAst(page *page.Page, parsedAst nodes.AbstractSyntaxTree) {
	// Second pass reactive content
	for _, node := range parsedAst {
		nodeContent := node.Content(page, rootAstNode)
		page.SetContent(nodeContent.HtmlContent.Content)
		page.AddStyle(nodeContent.CssContent)
		page.AddScript(nodeContent.JsContent)
		page.AddTwoScript(nodeContent.TwoScriptContent)

		recurseAst(page, node.Children())
	}
}
