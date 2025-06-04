package models

import (
	"hudson-newey/2web/src/compiler/lexer"
	"strings"
)

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

	// If the reducer for this event contains the referenced variable
	// e.g. $value = $value + 1 would have the reducer $value + 1
	// where the reducer contains the referenced variable
	//
	// Then we can set the variables "Reactive" property indicating that the
	// variable needs to use the least efficient updating method.
	//
	// A variable is considered reactive if its previous value is used in a
	// reducer.
	//
	// TODO: this event method should not be modifying the ReactiveVariable
	// properties
	varReferences := strings.Count(model.Reducer, variable.Name)
	if varReferences > 0 {
		variable.Reactive = true
	}
}
