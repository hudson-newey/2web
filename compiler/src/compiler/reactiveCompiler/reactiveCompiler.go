package reactiveCompiler

import (
	"hudson-newey/2web/src/models"
)

func CompileReactivity(
	filePath string,
	content string,
	varNodes []*models.ReactiveVariable,
) string {
	for _, varNode := range varNodes {
		// Ideally, slower reactive types would only target properties and events
		// that are effected.
		// Therefore, the fully "models.Reactive" variable reactivity class is a
		// superset of "models.Assignment" because some references to the variable
		// might not be fully reactive.
		//
		// This also means that each reactivity type should also make their own
		// element selectors so that subset of elements that abide by the reactivity
		// class are updated.
		//
		// TODO: I might be able to combine selectors for the same element that has
		// different property targets.
		if varNode.Type() >= models.Reactive {
			content = compileReactiveVar(content, varNode)
		} else if varNode.Type() >= models.Assignment {
			// TODO: explore if reactive and assignment reactivity are mutually
			// exclusive for variables, events, or props
			content = compileAssignmentVar(content, varNode)
		}

		// static props differ from truly static variables because static props
		// need runtime code to set the initial value of the prop
		//
		// e.g. <my-custom-element *value="$value"></my-custom-element>
		//
		// Most elements that have writable properties, also have an associated
		// attribute. So static properties only really apply to poorly designed
		// custom elements (e.g. web components).
		//
		// If the star is removed from this attribute, then it will become static
		// content that can be evaluated at runtime.
		// It is therefore recommended to remove the star from attributes if you
		// do not need to
		//
		// e.g. <input type="range" value="$value"></input>
		if varNode.Type() >= models.StaticProperty {
			content = compileStaticPropVar(content, varNode)
		}

		content = compileStatic(content, varNode)
	}

	return content
}
