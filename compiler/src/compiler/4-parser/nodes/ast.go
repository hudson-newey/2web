package nodes

import "github.com/hudson-newey/2web/_shared/lists"

type AbstractSyntaxTree []Node

func (ast AbstractSyntaxTree) reactiveVariables() []*reactiveVariableNode {
	return lists.FilterTypes[*reactiveVariableNode](ast)
}

func (ast AbstractSyntaxTree) reactiveProperties() []*reactivePropertyNode {
	return lists.FilterTypes[*reactivePropertyNode](ast)
}

func (ast AbstractSyntaxTree) reactiveEvents() []*reactiveEventNode {
	return lists.FilterTypes[*reactiveEventNode](ast)
}
