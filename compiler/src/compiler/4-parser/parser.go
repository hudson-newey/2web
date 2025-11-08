package parser

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/compiler/4-parser/grammar"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

func CreateAst(lexNodes []*lexer.V2LexNode) ast.AbstractSyntaxTree {
	var ast ast.AbstractSyntaxTree
	skipCount := 0

	for i := range lexNodes {
		if skipCount > 0 {
			skipCount--
			continue
		}

		nextNode, skipNext := processNode(i, lexNodes)
		skipCount = skipNext

		ast = append(ast, nextNode)
	}

	return ast
}

// TODO: Because the parser is VERY basic at the moment, we default to emitting
// a text node so that the old compiler can still process the text through the
// old string manipulation methods.
func processNode(index int, lexNodes []*lexer.V2LexNode) (ast.Node, int) {
	grammars := grammar.Rules
	for _, rule := range grammars {
		// Searches ahead of the current position in the lexed tokens to see if the
		// grammar definition matches.
		// We search ahead so that if the grammar definition is multiple tokens long,
		// we don't default to the text lexer state for leading tokens.
		peekBufferEnd := index + len(rule.Def)
		if peekBufferEnd > len(lexNodes) {
			continue
		}

		peekBuffer := lexNodes[index:peekBufferEnd]

		if rule.Matches(peekBuffer) {
			newNode := rule.Constructor(peekBuffer)
			return *newNode, len(rule.Def) - 1
		}
	}

	currentNode := lexNodes[index]
	return nodes.NewMarkupTextNode([]*lexer.V2LexNode{currentNode}), 0
}
