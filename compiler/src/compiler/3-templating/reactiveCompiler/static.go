package reactiveCompiler

import (
	"hudson-newey/2web/src/models"
	"strings"
)

func compileStatic(content string, varNode *models.ReactiveVariable) string {
	initialPropValue := strings.TrimPrefix(varNode.InitialValue, "\"")
	initialPropValue = strings.TrimPrefix(initialPropValue, "'")
	initialPropValue = strings.TrimSuffix(initialPropValue, "\"")
	initialPropValue = strings.TrimSuffix(initialPropValue, "'")

	content = strings.ReplaceAll(content, varNode.Name, initialPropValue)

	return content
}
