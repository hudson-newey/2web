package imports

import (
	"fmt"
	"hudson-newey/2web/src/models"
	"strings"
)

func ExpandComponentImports(
	workingPath string,
	content string,
	components []*models.Component,
) string {
	result := content

	for _, component := range components {
		result = expandImport(workingPath, content, component)
	}

	return result
}

func expandImport(workingPath string, content string, component *models.Component) string {
	templateSelector := fmt.Sprintf("<%s />", component.DomSelector)
	componentTemplate := component.HtmlContent(workingPath)

	result := strings.ReplaceAll(content, templateSelector, componentTemplate)

	return result
}
