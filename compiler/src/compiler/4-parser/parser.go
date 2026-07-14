package parser

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/compiler/4-parser/grammar"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

func CreateAst(lexNodes []*lexer.V2LexNode, grammars []grammar.Grammar) ast.AbstractSyntaxTree {
	var ast ast.AbstractSyntaxTree
	skipCount := 0

	for i := range lexNodes {
		if skipCount > 0 {
			skipCount--
			continue
		}

		nextNode, skipNext := processNode(i, lexNodes, grammars)
		skipCount = skipNext

		ast = append(ast, nextNode)
	}

	return ast
}

// TODO: Because the parser is VERY basic at the moment, we default to emitting
// a text node so that the old compiler can still process the text through the
// old string manipulation methods.
func processNode(index int, lexNodes []*lexer.V2LexNode, grammars []grammar.Grammar) (ast.Node, int) {
	for _, rule := range grammars {
		// Searches ahead of the current position in the lexed tokens to see if the
		// grammar definition matches.
		// We search ahead so that if the grammar definition is multiple tokens long,
		// we don't default to the text lexer state for leading tokens.
		peekBufferEnd := index + len(rule.Def)
		if peekBufferEnd > len(lexNodes) {
			continue
		}

		// TODO: Instead of passing the entire file into the Matches()
		// we should allow the Matches() call to progressively pull in tokens
		// peekBuffer := lexNodes[index:peekBufferEnd]
		peekBuffer := lexNodes[index:]

		matchedSubset := rule.Match(peekBuffer)
		if len(matchedSubset) > 0 {
			newNode := *rule.Constructor(peekBuffer)

			// Recurisvely match the ChildDef's
			childAst := CreateAst(matchedSubset, rule.ChildDefs)
			for _, child := range childAst {
				newNode.AddChild(child)
			}

			return newNode, len(rule.Def) - 1
		}
	}

	currentNode := lexNodes[index]
	return nodes.NewMarkupTextNode([]*lexer.V2LexNode{currentNode}), 0
}
