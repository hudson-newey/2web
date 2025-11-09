package templating

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/models/component"
)

// A function that takes an array of import nodes, evaluates, and returns the
// final file content with all of the imports expanded, inlined, and evaluated.
// This function takes an array of lexer nodes so that adding different import
// types e.g. component, css, JavaScript, html, etc... is the responsibility of
// this function instead of the compiler.
func expandImports(
	filePath string,
	pageModel *page.Page,
	importNodes []lexer.LexNode[lexer.ImportNode],
) {
	componentImports := []*models.Component{}

	for _, importNode := range importNodes {
		componentModel, err := component.FromNode(importNode, filePath)
		if err != nil {
			documentErrors.AddErrors(
				models.NewError(err.Error(), filePath, lexer.Position{}),
			)
			continue
		}

		componentImports = append(componentImports, &componentModel)
	}

	expandComponentImports(filePath, pageModel, componentImports)
}
