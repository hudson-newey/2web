package grammar

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

func wrapConstructor[T nodes.Node](ctor func([]*lexer.LexNode) T) func([]*lexer.LexNode) *nodes.Node {
	return func(lexNodes []*lexer.LexNode) *nodes.Node {
		newNode := ctor(lexNodes)
		var astNode nodes.Node = newNode
		return &astNode
	}
}
