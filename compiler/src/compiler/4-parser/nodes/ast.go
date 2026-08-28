package nodes

import "github.com/hudson-newey/2web/_shared/lists"

type AbstractSyntaxTree []Node

func (ast AbstractSyntaxTree) reactiveVariables() []*reactiveVariableNode {
	matches := []*reactiveVariableNode{}
	for _, n := range ast {
		matches = append(matches, n.Children().reactiveVariables()...)
	}

	matches = append(matches, lists.FilterTypes[*reactiveVariableNode](ast)...)
	return matches
}

func (ast AbstractSyntaxTree) reactiveProperties() []*reactivePropertyNode {
	matches := []*reactivePropertyNode{}
	for _, n := range ast {
		matches = append(matches, n.Children().reactiveProperties()...)
	}

	matches = append(matches, lists.FilterTypes[*reactivePropertyNode](ast)...)
	return matches
}

func (ast AbstractSyntaxTree) reactiveEvents() []*reactiveEventNode {
	matches := []*reactiveEventNode{}
	for _, n := range ast {
		matches = append(matches, n.Children().reactiveEvents()...)
	}

	matches = append(matches, lists.FilterTypes[*reactiveEventNode](ast)...)
	return matches
}
