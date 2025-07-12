package models

import lexer "hudson-newey/2web/src/compiler/2-lexer"

type ReactiveProperty struct {
	PropName string
	VarName  string
	Variable *ReactiveVariable
	Node     *lexer.LexNode[lexer.PropNode]
}

func (model *ReactiveProperty) BindVariable(variable *ReactiveVariable) {
	if model.Variable != nil {
		panic("Binding already exists")
	} else if variable.Name != model.VarName {
		panic("Attempted to bind variable to mismatched property")
	}

	model.Variable = variable
}
