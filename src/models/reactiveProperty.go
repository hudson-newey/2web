package models

import "hudson-newey/2web/src/lexer"

type ReactiveProperty struct {
	Node     *lexer.LexNode[lexer.PropNode]
	PropName string
	VarName  string
	Variable *ReactiveVariable
}

func (model *ReactiveProperty) BindVariable(variable *ReactiveVariable) {
	if model.Variable != nil {
		panic("Binding already exists")
	} else if variable.Name != model.VarName {
		panic("Attempted to bind variable to mismatched property")
	}

	model.Variable = variable
}
