package parser

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/compiler/4-parser/grammar"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

// A buffer of nodes being processed from the parser.
// Once a grammar rule matches the nodes in the buffer, they are consumed,
// the grammar's constructor is called to create a new AST node, and the buffer
// is cleared for the next set of nodes.
var nodeBuffer = []*lexer.V2LexNode{}

func CreateAst(lexNode []lexer.V2LexNode) ast.AbstractSyntaxTree {
	var ast ast.AbstractSyntaxTree
	for _, node := range lexNode {
		ast = append(ast, processNode(&node))
	}

	return ast
}

// TODO: Because the parser is VERY basic at the moment, we default to emitting
// a text node so that the old compiler can still process the text through the
// old string manipulation methods.
func processNode(lexNode *lexer.V2LexNode) ast.Node {
	nodeBuffer = append(nodeBuffer, lexNode)

	grammars := grammar.Rules
	for _, rule := range grammars {
		if rule.Matches(nodeBuffer) {
			newNode := rule.Constructor(nodeBuffer)

			nodeBuffer = []*lexer.V2LexNode{}

			return *newNode
		}
	}

	return nodes.NewMarkupTextNode([]*lexer.V2LexNode{lexNode})
}
