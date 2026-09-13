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

	// dependencies maps a variable's selector to the selectors of the
	// reactive variables referenced in its initial value expression (its
	// "computation" dependencies). Variables without references are absent.
	dependencies map[string][]string

	// runtimeVariables maps a variable's selector to the runtime JavaScript
	// variable name allocated for it during compilation. Only variables that
	// need a runtime representation (event assigned variables, and variables
	// that other variables are computed from) get one.
	runtimeVariables map[string]string
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
		Variables:        variables,
		dependencies:     map[string][]string{},
		runtimeVariables: map[string]string{},
	}

	// Resolve the variable-to-variable dependency edges: a variable whose
	// initial value references another reactive variable is computed from it.
	for _, variable := range variables {
		deps := []string{}
		for _, candidate := range variables {
			if candidate == variable {
				continue
			}

			if strings.Contains(variable.initialValue, candidate.selector()) {
				deps = append(deps, candidate.selector())
			}
		}

		if len(deps) > 0 {
			index.dependencies[variable.selector()] = deps
		}
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

// ---------------------------------------------------------------------------
// Dependent (computed) variables
// ---------------------------------------------------------------------------

// IsDerived returns whether the variable is computed from other reactive
// variables (its initial value references at least one other reactive
// variable).
func (index *ReactiveIndex) IsDerived(variable *reactiveVariableNode) bool {
	return len(index.DependenciesOf(variable)) > 0
}

// DependenciesOf returns the selectors of the reactive variables that the
// variable is computed from.
func (index *ReactiveIndex) DependenciesOf(variable *reactiveVariableNode) []string {
	return index.dependencies[variable.selector()]
}

// RegisterRuntimeVariable records the runtime JavaScript variable name that
// was allocated for a reactive variable.
func (index *ReactiveIndex) RegisterRuntimeVariable(selector string, runtimeName string) {
	index.runtimeVariables[selector] = runtimeName
}

// RuntimeVariableName returns the runtime JavaScript variable name allocated
// for a reactive variable, or empty when the variable has no runtime
// representation.
func (index *ReactiveIndex) RuntimeVariableName(selector string) string {
	return index.runtimeVariables[selector]
}

// IsRuntime returns whether the variable is directly assigned by an event,
// which means it needs a runtime representation.
func (index *ReactiveIndex) IsRuntime(variable *reactiveVariableNode) bool {
	for _, event := range index.Events {
		if event.sink == variable.selector() {
			return true
		}
	}

	return false
}

// HasRuntimeDependency returns whether any variable in the variable's
// dependency closure has a runtime representation. When that is the case, the
// variable is recomputed by the runtime dependency's update cascade instead of
// being wired here.
func (index *ReactiveIndex) HasRuntimeDependency(variable *reactiveVariableNode) bool {
	for _, dependency := range index.dependencyClosure(variable, map[string]bool{}) {
		if index.IsRuntime(dependency) {
			return true
		}
	}

	return false
}

// HasDerivedDependents returns whether any other variable is computed from
// this variable. Such variables need this one to have a runtime
// representation so that they can be recomputed when it changes.
func (index *ReactiveIndex) HasDerivedDependents(variable *reactiveVariableNode) bool {
	selector := variable.selector()

	for _, dependencies := range index.dependencies {
		if slices.Contains(dependencies, selector) {
			return true
		}
	}

	return false
}

// IsUnused returns whether the variable has no reactive properties, no
// events, and isn't referenced by any other reactive variable.
func (index *ReactiveIndex) IsUnused(variable *reactiveVariableNode) bool {
	if len(index.FindDependentProperties(variable)) > 0 {
		return false
	}

	if len(index.FindDependentEvents(variable)) > 0 {
		return false
	}

	return !index.HasDerivedDependents(variable)
}

// dependencyClosure returns every variable that the given variable transitively
// depends on (not including the variable itself).
func (index *ReactiveIndex) dependencyClosure(variable *reactiveVariableNode, visited map[string]bool) []*reactiveVariableNode {
	closure := []*reactiveVariableNode{}

	for _, dependencySelector := range index.DependenciesOf(variable) {
		if visited[dependencySelector] {
			continue
		}

		visited[dependencySelector] = true

		dependency := index.VariableBySelector(dependencySelector)
		if dependency == nil {
			continue
		}

		closure = append(closure, dependency)
		closure = append(closure, index.dependencyClosure(dependency, visited)...)
	}

	return closure
}

// ResolveExpression returns the JavaScript expression that evaluates the
// variable's value at runtime.
//
// Reactive variables referenced by the expression are substituted with the
// runtime variable names that were allocated for them (see
// RegisterRuntimeVariable). Static variables referenced by the expression are
// substituted with their own recursively resolved initial expression, so a
// chain of computed variables collapses into a single expression.
func (index *ReactiveIndex) ResolveExpression(variable *reactiveVariableNode) (string, bool) {
	return index.resolveExpression(variable, map[string]bool{})
}

func (index *ReactiveIndex) resolveExpression(variable *reactiveVariableNode, visiting map[string]bool) (string, bool) {
	selector := variable.selector()

	if visiting[selector] {
		return "", false
	}

	visiting[selector] = true
	defer delete(visiting, selector)

	expression := variable.initialValue

	for _, dependencySelector := range index.DependenciesOf(variable) {
		dependency := index.VariableBySelector(dependencySelector)
		if dependency == nil {
			continue
		}

		// Runtime dependencies are substituted with the runtime variable that
		// was allocated for them in the reactive runtime scope.
		if runtimeName := index.RuntimeVariableName(dependencySelector); runtimeName != "" {
			expression = replaceSelector(expression, dependencySelector, runtimeName)
			continue
		}

		// Static dependencies are substituted with their resolved initial
		// expression (recursively, for chains of computed variables).
		dependencyExpression, ok := index.resolveExpression(dependency, visiting)
		if !ok {
			return "", false
		}

		expression = replaceSelector(expression, dependencySelector, "("+dependencyExpression+")")
	}

	return expression, true
}

// replaceSelector replaces every occurrence of a variable selector with the
// given replacement, respecting identifier boundaries so that replacing
// "$count" doesn't corrupt a reference to "$count2".
func replaceSelector(source string, selector string, replacement string) string {
	var result strings.Builder

	for i := 0; i < len(source); {
		if source[i] == selector[0] && i+len(selector) <= len(source) {
			end := i + len(selector)

			before := byte(' ')
			if i > 0 {
				before = source[i-1]
			}

			after := byte(' ')
			if end < len(source) {
				after = source[end]
			}

			if source[i:end] == selector && !isIdentifierByte(before) && !isIdentifierByte(after) {
				result.WriteString(replacement)
				i = end
				continue
			}
		}

		result.WriteByte(source[i])
		i++
	}

	return result.String()
}

func isIdentifierByte(b byte) bool {
	return b == '_' || b == '$' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

// VariableBySelector returns the variable declared with the given selector,
// or nil when no such variable exists.
func (index *ReactiveIndex) VariableBySelector(selector string) *reactiveVariableNode {
	for _, variable := range index.Variables {
		if variable.selector() == selector {
			return variable
		}
	}

	return nil
}

// ComputedDependents returns every variable that is (transitively) computed
// from the given variable, ordered so that every variable appears after the
// variables it is computed from. A nil result means the dependency graph is
// cyclic and no ordering exists.
func (index *ReactiveIndex) ComputedDependents(variable *reactiveVariableNode) []*reactiveVariableNode {
	ordered := []*reactiveVariableNode{}
	visiting := map[string]bool{}
	done := map[string]bool{}

	// The closure is collected in post order (a variable after its own
	// dependents) and reversed so that the result is in dependency order (a
	// variable before the variables computed from it).
	var visit func(current *reactiveVariableNode) bool
	visit = func(current *reactiveVariableNode) bool {
		selector := current.selector()

		if done[selector] {
			return true
		}

		if visiting[selector] {
			// The dependency graph is cyclic.
			return false
		}

		visiting[selector] = true

		for _, candidate := range index.Variables {
			deps := index.dependencies[candidate.selector()]
			if !slices.Contains(deps, selector) {
				continue
			}

			if !visit(candidate) {
				return false
			}
		}

		delete(visiting, selector)
		done[selector] = true
		ordered = append(ordered, current)
		return true
	}

	if !visit(variable) {
		return nil
	}

	// The root itself is visited (as the entry point of the closure) but isn't
	// part of its own dependents.
	ordered = ordered[:len(ordered)-1]

	// Reverse the post order into dependency order.
	for i, j := 0, len(ordered)-1; i < j; i, j = i+1, j-1 {
		ordered[i], ordered[j] = ordered[j], ordered[i]
	}

	return ordered
}
