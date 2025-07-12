package component

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/models"
	"strings"
)

// An incrementing id that can be used to uniquely identify the component.
// Note that if used, it will be compressed into a base 32 DOM selector.
var nextNodeId uint64 = 0

func FromNode(node lexer.LexNode[lexer.ImportNode], importedFrom string) (models.Component, error) {
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
		DomSelector:  domSelector,
		Identifier:   nextNodeId,
		ImportPath:   importPath,
		ImportedFrom: importedFrom,
	}

	nextNodeId += 1

	return componentModel, nil
}
