package grammar

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

func wrapConstructor[T nodes.Node](ctor func([]*lexer.V2LexNode, *nodes.ParseContext) T) func([]*lexer.V2LexNode, *nodes.ParseContext) *nodes.Node {
	return func(lexNodes []*lexer.V2LexNode, context *nodes.ParseContext) *nodes.Node {
		newNode := ctor(lexNodes, context)
		var astNode nodes.Node = newNode
		return &astNode
	}
}
