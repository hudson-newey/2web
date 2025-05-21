package component

import (
	"fmt"
	"hudson-newey/2web/src/compiler/lexer"
	"hudson-newey/2web/src/models"
)

func FromNode(node lexer.LexNode[lexer.ImportNode]) (models.Component, error) {
	if len(node.Tokens) != 3 {
		errorMessage := fmt.Errorf("malformed import statement:\n\tExpected: import ComponentName from \"path/to/asset.html\";\n\tFound: %s", node.Selector)
		return models.Component{}, errorMessage
	}

	domSelector := node.Tokens[1]
	importPath := node.Tokens[2]

	componentModel := models.Component{
		DomSelector: domSelector,
		ImportPath:  importPath,
	}

	return componentModel, nil
}
