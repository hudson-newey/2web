package compiler

import (
	"hudson-newey/2web/src/models"
	"strings"
)

func compileReactivity(
	filePath string,
	content string,
	varNodes []models.ReactiveVariable,
	propNodes []models.ReactiveProperty,
) string {
	for _, variable := range varNodes {
		content = strings.ReplaceAll(content, mustacheStartToken+variable.Name+mustacheEndToken, variable.InitialValue)
		content = strings.ReplaceAll(content, variable.Name, variable.InitialValue)
	}

	return content
}
