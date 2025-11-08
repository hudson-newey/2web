package grammar

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"strings"
)

type definition []lexerTokens.LexToken

type Grammar struct {
	// A sequence of tokens that define the reactive variable
	Def definition

	// A constructor function to create a node from the tokens
	Constructor func(lexNodes []*lexer.V2LexNode) *ast.Node
}

func (model *Grammar) Matches(lexNodes []*lexer.V2LexNode) bool {
	if len(lexNodes) != len(model.Def) {
		return false
	}

	for i, token := range model.Def {
		if lexNodes[i].Token != token {
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
