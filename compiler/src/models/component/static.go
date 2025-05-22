package component

import (
	"fmt"
	"hudson-newey/2web/src/compiler/lexer"
	"hudson-newey/2web/src/models"
	"strings"
)

func FromNode(node lexer.LexNode[lexer.ImportNode]) (models.Component, error) {
	if len(node.Tokens) != 3 {
		errorMessage := fmt.Errorf("malformed import statement:\n\tExpected: import ComponentName from \"path/to/asset.html\";\n\tFound: %s", node.Selector)
		return models.Component{}, errorMessage
	}

	domSelector := node.Tokens[0]
	importPath := node.Tokens[2]

	// Because the import will contain leading and trailing quotation marks,
	// in it's token, I remove them so that we don't have to deal with them
	// later.
	importPath, _ = strings.CutSuffix(importPath, "\"")
	importPath, _ = strings.CutPrefix(importPath, "\"")

	componentModel := models.Component{
		DomSelector: domSelector,
		ImportPath:  importPath,
	}

	return componentModel, nil
}
