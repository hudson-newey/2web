package parser

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/grammar"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

func CreateAst(
	lexNodes []*lexer.V2LexNode,
	grammars []grammar.Grammar,
	withTextMatcher bool,
	context *nodes.ParseContext,
) nodes.AbstractSyntaxTree {
	var ast nodes.AbstractSyntaxTree
	skipCount := 0

	for i := range lexNodes {
		if skipCount > 0 {
			skipCount--
			continue
		}

		nextNode, skipNext := processNode(i, lexNodes, grammars, withTextMatcher, context)
		skipCount = skipNext

		if nextNode != nil {
			ast = append(ast, nextNode)
		}
	}

	return ast
}

// TODO: Because the parser is VERY basic at the moment, we default to emitting
// a text node so that the old compiler can still process the text through the
// old string manipulation methods.
func processNode(
	index int,
	lexNodes []*lexer.V2LexNode,
	grammars []grammar.Grammar,
	withTextMatcher bool,
	context *nodes.ParseContext,
) (nodes.Node, int) {
	for _, rule := range grammars {
		// Searches ahead of the current position in the lexed tokens to see if the
		// grammar definition matches.
		// We search ahead so that if the grammar definition is multiple tokens long,
		// we don't default to the text lexer state for leading tokens.
		//
		// Optional tokens in the definition can be missing from the input, so the
		// look-ahead only needs to cover the required tokens.
		peekBufferEnd := index + rule.MinimumTokenCount()
		if peekBufferEnd > len(lexNodes) {
			continue
		}

		// TODO: Instead of passing the entire file into the Matches()
		// we should allow the Matches() call to progressively pull in tokens
		// peekBuffer := lexNodes[index:peekBufferEnd]
		peekBuffer := lexNodes[index:]

		matchedSubset := rule.Match(peekBuffer, context)
		if len(matchedSubset) > 0 {
			// The node constructors receive the matched subset (rather than the
			// whole tail) so that degraded nodes can render exactly the source
			// that the grammar consumed when reporting syntax errors.
			newNode := *rule.Constructor(matchedSubset, context)

			// Recurisvely match the ChildDef's
			childAst := CreateAst(matchedSubset, rule.ChildDefs, false, context)
			for _, child := range childAst {
				newNode.AddChild(child)
			}

			return newNode, len(matchedSubset) - 1
		}
	}

	if withTextMatcher {
		currentNode := lexNodes[index]
		return nodes.NewMarkupTextNode([]*lexer.V2LexNode{currentNode}), 0
	}
	return nil, 0
}
