package grammar

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"strings"
)

type grammar struct {
	// A sequence of tokens that define the reactive variable
	Def definition

	// A constructor function to create a node from the tokens
	Constructor func(lexNodes []*lexer.V2LexNode) *ast.Node
}

func (model *grammar) Matches(lexNodes []*lexer.V2LexNode) bool {
	// If we have not processed enough tokens to have a match yet, we can quickly
	// return false.
	if len(lexNodes) < len(model.Def) {
		return false
	}

	// Because a grammar definition can sometimes come after multiple non-matching
	// tokens, we want to check only the last N tokens where N is the length of
	// the grammar definition.
	matchingSubset := lexNodes[len(lexNodes)-len(model.Def):]

	for i, token := range model.Def {
		if matchingSubset[i].Token != token {
			return false
		}
	}

	return true
}

func (model *definition) String() string {
	var tokens []string
	for _, token := range *model {
		tokens = append(tokens, token.String())
	}
	return "[" + strings.Join(tokens, ", ") + "]"
}
