package parser

import lexer "hudson-newey/2web/src/compiler/2-lexer"

type AbstractSyntaxTree []Node

func CreateAst(structure []lexer.V2LexNode) AbstractSyntaxTree {
	return AbstractSyntaxTree{}
}
