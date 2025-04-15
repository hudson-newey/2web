package reactiveCompiler

import (
	"hudson-newey/2web/src/models"
	"strings"
)

func compileStatic(content string, varNode *models.ReactiveVariable) string {
	// TODO: remove this hard coded hack
	content = strings.ReplaceAll(
		content,
		"{{ "+varNode.Name+" }}",
		varNode.InitialValue,
	)
	content = strings.ReplaceAll(content, varNode.Name, varNode.InitialValue)

	return content
}
