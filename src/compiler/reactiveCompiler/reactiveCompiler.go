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
		if varNode.Type() >= models.Reactive {
			content = compileReactiveTemplate(content, varNode)
			continue
		}

		// TODO: this should have a "continue". But the implementation is bugged
		// and relies on the static compilation to do variable substitution.
		//
		// It works for the moment and I just want to get a working version, so I'm
		// leaving a comment for now.
		if varNode.Type() >= models.Assignment {
			content = compileAssignmentProp(content, varNode)
		}

		if varNode.Type() >= models.StaticProperty {
			content = compileStaticProperty(content, varNode)
		}

		content = compileStatic(content, varNode)
	}

	return content
}
