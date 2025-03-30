package models

import "hudson-newey/2web/src/lexer"

type ReactiveEvent struct {
	EventName string
	VarName   string
	Reducer   string
	Variable  *ReactiveVariable
	Node      *lexer.LexNode[lexer.EventNode]
}

func (model *ReactiveEvent) BindVariable(variable *ReactiveVariable) {
	if variable.Name != model.VarName {
		panic("Attempted to bind variable to mismatched property")
	}

	model.Variable = variable
}
