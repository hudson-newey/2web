package models

import (
	"hudson-newey/2web/src/lexer"
)

type ReactiveProperty struct {
	Name    string
	Binding string
}

func (model *ReactiveProperty) FromNode(node lexer.LexNode[lexer.PropNode]) {
	model.Name = node.Tokens[0]
	model.Binding = node.Tokens[1]
}
