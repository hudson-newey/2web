package nodes

import (
	"slices"
	"sort"
	"strings"
)

// ReactiveIndex is a pre-resolved view of a page's reactive graph.
//
// Resolving the dependency relationships between reactive variables, properties,
// and events used to be done by walking the entire page AST inside the
// dependency filters. That made reactivity compilation O(variables x properties
// x nodes): every variable re-walked the whole tree for every candidate
// property and event.
//
// The index is built once per page (from the page's root AST) and every
// dependency query is answered from it.
type ReactiveIndex struct {
	Variables  []*reactiveVariableNode
	Properties []indexedProperty
	Events     []indexedEvent
}

type indexedProperty struct {
	node *reactivePropertyNode
	// The selectors of the reactive variables that the property's reducer
	// references.
	dependencies []string
}

type indexedEvent struct {
	node *reactiveEventNode
	// The selector of the variable that the event assigns to, or empty when
	// the assignment target isn't a declared reactive variable.
	sink string
	// The selectors of the reactive variables that the event reducer reads.
	dependencies []string
}

// BuildReactiveIndex walks the AST once and resolves every reactive
// relationship up front.
func BuildReactiveIndex(ast AbstractSyntaxTree) *ReactiveIndex {
	variables := ast.reactiveVariables()

	index := &ReactiveIndex{
		Variables: variables,
	}

	for _, property := range ast.reactiveProperties() {
		index.Properties = append(index.Properties, indexedProperty{
			node:         property,
			dependencies: variableSelectorsReferencedBy(property.reducer, variables),
		})
	}

	for _, event := range ast.reactiveEvents() {
		index.Events = append(index.Events, indexedEvent{
			node:         event,
			sink:         eventSinkSelector(event.assignmentSink, variables),
			dependencies: variableSelectorsReferencedBy(event.assignmentExpr, variables),
		})
	}

	sort.Slice(index.Variables, func(i, j int) bool {
		return index.Variables[i].variableName < index.Variables[j].variableName
	})

	return index
}

// FindDependentProperties returns every property whose reducer references the
// given variable.
func (index *ReactiveIndex) FindDependentProperties(variable *reactiveVariableNode) []*reactivePropertyNode {
	selector := variable.selector()

	matches := []*reactivePropertyNode{}
	for _, property := range index.Properties {
		if slices.Contains(property.dependencies, selector) {
			matches = append(matches, property.node)
		}
	}

	return matches
}

// FindDependentEvents returns every event that reads or assigns to the given
// variable.
func (index *ReactiveIndex) FindDependentEvents(variable *reactiveVariableNode) []*reactiveEventNode {
	selector := variable.selector()

	matches := []*reactiveEventNode{}
	for _, event := range index.Events {
		if slices.Contains(event.dependencies, selector) || event.sink == selector {
			matches = append(matches, event.node)
		}
	}

	return matches
}

// variableSelectorsReferencedBy returns the selectors of every variable whose
// selector appears in the given source fragment.
func variableSelectorsReferencedBy(source string, variables []*reactiveVariableNode) []string {
	selectors := []string{}
	for _, variable := range variables {
		if containsSelector(source, variable.selector()) {
			selectors = append(selectors, variable.selector())
		}
	}

	return selectors
}

// eventSinkSelector returns the selector of the variable with the given name,
// or empty when no declared variable matches.
func eventSinkSelector(sinkName string, variables []*reactiveVariableNode) string {
	for _, variable := range variables {
		if variable.selector() == sinkName {
			return variable.selector()
		}
	}

	return ""
}

// containsSelector checks whether a reducer references a variable selector.
// This mirrors the contains checks that the reactive nodes used to perform
// inline during compilation.
func containsSelector(source string, selector string) bool {
	return strings.Contains(source, selector)
}
