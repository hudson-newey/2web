package templating

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/compiler/5-templating/reactiveCompiler"
	"hudson-newey/2web/src/content/assets"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/content/xslt"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/models/reactiveEvent"
	"hudson-newey/2web/src/models/reactiveProperty"
	"hudson-newey/2web/src/models/reactiveVariable"
	"os"
	"strings"
)

func Compile(filePath string, parsedAst ast.AbstractSyntaxTree) page.Page {
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

	recurseAst(&pageModel, parsedAst)
	addRouteAssets(&pageModel)

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
	pageModel.Html.Content = expandTextNodes(pageModel.Html.Content)

	reactiveVariables := []*models.ReactiveVariable{}

	propertyNodes := lexer.FindPropNodes[lexer.PropNode](pageModel.Html.Content, propertyPrefix)
	eventNodes := lexer.FindPropNodes[lexer.EventNode](pageModel.Html.Content, eventPrefix)

	for _, node := range pageModel.TwoScript {
		lineCommentNodes := lexer.FindNodes[lexer.LineCommentNode](node.Content, lineCommentStartToken, newLineToken)
		lineCommentRemoved := node.Content
		for _, commentNode := range lineCommentNodes {
			lineCommentRemoved = strings.ReplaceAll(lineCommentRemoved, commentNode.Selector, "")
		}

		// We cannot do this the same time as line comments because someone might
		// decide to put line comments inside of a block comment (or vice versa).
		// Meaning that the selectors would no longer match.
		blockCommentNodes := lexer.FindNodes[lexer.BlockCommentNode](lineCommentRemoved, blockCommentStartToken, blockCommentEndToken)
		noCommentsContent := lineCommentRemoved
		for _, commentNode := range blockCommentNodes {
			noCommentsContent = strings.ReplaceAll(noCommentsContent, commentNode.Selector, "")
		}

		variableNodes := lexer.FindNodes[lexer.VarNode](noCommentsContent, variableToken, statementEndToken)

		for _, variableNode := range variableNodes {
			variableModel, err := reactiveVariable.FromNode(variableNode)
			if err != nil {
				varError := models.NewError(
					err.Error(),
					filePath,
					lexer.Position{},
				)

				pageModel.Errors.AddError(&varError)
				continue
			}

			reactiveVariables = append(reactiveVariables, &variableModel)
		}
	}

	for _, node := range propertyNodes {
		property, err := reactiveProperty.FromNode(node)
		if err != nil {
			varError := models.NewError(
				err.Error(),
				filePath,
				lexer.Position{},
			)

			pageModel.Errors.AddError(&varError)
			continue
		}

		// find the reactive variable that matches the binding name
		foundAssociatedProperty := false
		for _, variable := range reactiveVariables {
			if variable.Name == property.VarName {
				variable.AddProp(&property)
				property.BindVariable(variable)
				foundAssociatedProperty = true
				break
			}
		}

		if !foundAssociatedProperty {
			errorMessage := fmt.Sprintf("could not find compiler variable '%s' for property %s", property.VarName, property.Node.Selector)
			varError := models.NewError(
				errorMessage,
				filePath,
				lexer.Position{},
			)

			pageModel.Errors.AddError(&varError)
			continue
		}
	}

	for _, node := range eventNodes {
		event, err := reactiveEvent.FromNode(node)
		if err != nil {
			varError := models.NewError(
				err.Error(),
				filePath,
				lexer.Position{},
			)

			pageModel.Errors.AddError(&varError)
			continue
		}

		foundAssociatedProperty := false
		for _, variable := range reactiveVariables {
			if variable.Name == event.VarName {
				variable.AddEvent(&event)
				event.BindVariable(variable)
				foundAssociatedProperty = true
				break
			}
		}

		if !foundAssociatedProperty {
			errorMessage := fmt.Sprintf("could not find compiler variable '%s' for event %s", event.VarName, event.Node.Selector)
			varError := models.NewError(
				errorMessage,
				filePath,
				lexer.Position{},
			)

			pageModel.Errors.AddError(&varError)
			continue
		}
	}

	reactiveCompiler.CompileReactivity(filePath, &pageModel, reactiveVariables)

	return pageModel
}

func recurseAst(page *page.Page, parsedAst ast.AbstractSyntaxTree) {
	// First pass to establish page content
	// We need this so that when we get to ast nodes that replace page content,
	// it has the full page context.
	for _, node := range parsedAst {
		if node.Type() != "MarkupTextNode" {
			continue
		}

		nodeContent := node.Content(page)
		page.SetContent(nodeContent.HtmlContent)
		page.AddStyle(nodeContent.CssContent)
		page.AddScript(nodeContent.JsContent)
		page.AddTwoScript(nodeContent.TwoScriptContent)

		recurseAst(page, node.Children())
	}

	// Second pass reactive content
	for _, node := range parsedAst {
		if node.Type() == "MarkupTextNode" {
			continue
		}

		nodeContent := node.Content(page)
		page.SetContent(nodeContent.HtmlContent)
		page.AddStyle(nodeContent.CssContent)
		page.AddScript(nodeContent.JsContent)
		page.AddTwoScript(nodeContent.TwoScriptContent)

		recurseAst(page, node.Children())
	}
}
