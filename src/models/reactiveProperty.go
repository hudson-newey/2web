package models

import "hudson-newey/2web/src/lexer"

type ReactiveProperty struct {
	PropName    string
	Node        *lexer.LexNode[lexer.PropNode]
	BindingName string
	Variable    *ReactiveVariable
}

func (model *ReactiveProperty) AddBinding(variable *ReactiveVariable) {
	if model.Variable != nil {
		panic("Binding already exists")
	} else if variable.Name != model.BindingName {
		panic("Attempted to bind variable to mismatched property")
	}

	model.Variable = variable
}
