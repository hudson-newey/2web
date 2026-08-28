package nodes

import "github.com/hudson-newey/2web/_shared/lists"

type AbstractSyntaxTree []Node

func (ast AbstractSyntaxTree) reactiveVariables() []*reactiveVariableNode {
	matches := []*reactiveVariableNode{}
	for _, n := range ast {
		matches = append(matches, n.Children().reactiveVariables()...)
	}

	return lists.FilterTypes[*reactiveVariableNode](ast)
}

func (ast AbstractSyntaxTree) reactiveProperties() []*reactivePropertyNode {
	matches := []*reactivePropertyNode{}
	for _, n := range ast {
		matches = append(matches, n.Children().reactiveProperties()...)
	}

	return lists.FilterTypes[*reactivePropertyNode](ast)
}

func (ast AbstractSyntaxTree) reactiveEvents() []*reactiveEventNode {
	matches := []*reactiveEventNode{}
	for _, n := range ast {
		matches = append(matches, n.Children().reactiveEvents()...)
	}

	return lists.FilterTypes[*reactiveEventNode](ast)
}
