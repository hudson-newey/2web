package nodes

import (
	"sort"

	"hudson-newey/2web/src/debugger"
)

// CollectDebugInfo walks the AST and records the reactive graph that the
// compiler built for a page: every reactive variable (with its reactivity
// class), every reactive property (with the variables its reducer depends on),
// and every reactive event (with its sink and dependencies).
//
// This powers the `__2web.debug.json` file that the browser devtools extension
// renders.
func CollectDebugInfo(inputPath string, ast AbstractSyntaxTree) debugger.PageDebugInfo {
	pageDebug := debugger.PageDebugInfo{
		Page:       inputPath,
		Variables:  []debugger.VariableInfo{},
		Properties: []debugger.PropertyInfo{},
		Events:     []debugger.EventInfo{},
	}

	for _, variable := range ast.reactiveVariables() {
		pageDebug.Variables = append(pageDebug.Variables, debugger.VariableInfo{
			Name:            variable.selector(),
			InitialValue:    variable.initialValue,
			ReactivityClass: reactivityLevelName(variable.reactivityLevel(ast)),
		})
	}

	for _, property := range ast.reactiveProperties() {
		pageDebug.Properties = append(pageDebug.Properties, debugger.PropertyInfo{
			PropName:     property.propName,
			Reducer:      property.reducer,
			Dependencies: dependencyNames(property.reactiveVariableDeps(ast)),
		})
	}

	for _, event := range ast.reactiveEvents() {
		pageDebug.Events = append(pageDebug.Events, debugger.EventInfo{
			EventName:    event.eventName,
			Reducer:      event.reducer,
			Sink:         event.sinkName(ast),
			Dependencies: dependencyNames(event.reactiveVariableDeps(ast)),
		})
	}

	sortVariables(pageDebug.Variables)
	sortProperties(pageDebug.Properties)
	sortEvents(pageDebug.Events)

	return pageDebug
}

func dependencyNames(variables []*reactiveVariableNode) []string {
	names := []string{}
	for _, variable := range variables {
		names = append(names, variable.selector())
	}

	sort.Strings(names)
	return names
}

// sinkName resolves the variable that an event assigns to without panicking.
// The compiler's reactiveVariableSink panics when the assignment target isn't a
// declared reactive variable, which is surfaced as a compiler error elsewhere.
func (m *reactiveEventNode) sinkName(ast AbstractSyntaxTree) string {
	sink := ""
	func() {
		defer func() {
			if recover() != nil {
				sink = ""
			}
		}()

		if matched := m.reactiveVariableSink(ast); matched != nil {
			sink = matched.selector()
		}
	}()

	return sink
}

func reactivityLevelName(level reactivityLevel) string {
	switch level {
	case unused:
		return "unused"
	case static:
		return "static"
	case staticProperty:
		return "static-property"
	case assignment:
		return "assignment"
	case reactive:
		return "reactive"
	}

	return "unknown"
}

func sortVariables(variables []debugger.VariableInfo) {
	sort.Slice(variables, func(i, j int) bool {
		return variables[i].Name < variables[j].Name
	})
}

func sortProperties(properties []debugger.PropertyInfo) {
	sort.Slice(properties, func(i, j int) bool {
		return properties[i].PropName < properties[j].PropName
	})
}

func sortEvents(events []debugger.EventInfo) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].EventName < events[j].EventName
	})
}
