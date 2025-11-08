package grammar

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
)

func wrapConstructor[T ast.Node](ctor func([]*lexer.V2LexNode) T) func([]*lexer.V2LexNode) *ast.Node {
	return func(lexNodes []*lexer.V2LexNode) *ast.Node {
		newNode := ctor(lexNodes)
		var astNode ast.Node = newNode
		return &astNode
	}
}
