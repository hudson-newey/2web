package templating

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
	optimizer "hudson-newey/2web/src/compiler/6-optimizer"
	"hudson-newey/2web/src/content/assets"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/content/xslt"
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
		pageModel.SetContent(
			html.FromContent(string(fileContent)),
		)
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

	// We want to exclude Markdown, xml, and xslt files from this processing step
	// because our element ref symbol is a hashtag which has meaning in these
	// formats.
	//
	// TODO: We should add support for element refs in Markdown, xml and xslt
	// files.
	if assets.IsMarkupFile(filePath) &&
		!markdown.IsMarkdownFile(filePath) &&
		!xslt.IsXsltFile(filePath) {
		pageModel.Html.Content = expandElementRefs(pageModel.Html.Content)
	}

	// "Mustache like" expressions e.g. {{ $count }} are a shorthand for an
	// element with only innerText.
	// Therefore, we expand all of the mustache expressions before finding
	// reactive property tokens
	// pageModel.Html.Content = expandTextNodes(pageModel.Html.Content)

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
		page.SetContent(
			html.FromContent(page.Html.Content + markupContent),
		)

		recurseAstMarkup(page, node.Children())
	}
}

func recurseAst(page *page.Page, parsedAst nodes.AbstractSyntaxTree) {
	// Second pass reactive content
	for _, node := range parsedAst {
		nodeContent := node.Content(page, rootAstNode)
		page.SetContent(nodeContent.HtmlContent)
		page.AddStyle(nodeContent.CssContent)
		page.AddScript(nodeContent.JsContent)
		page.AddTwoScript(nodeContent.TwoScriptContent)

		recurseAst(page, node.Children())
	}
}
