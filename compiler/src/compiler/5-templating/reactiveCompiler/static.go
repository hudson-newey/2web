package reactiveCompiler

import (
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileStatic(
	pageModel *page.Page,
	varNode *models.ReactiveVariable,
) {
	content := pageModel.Html.Content

	initialPropValue := strings.TrimPrefix(varNode.InitialValue, "\"")
	initialPropValue = strings.TrimPrefix(initialPropValue, "'")
	initialPropValue = strings.TrimSuffix(initialPropValue, "\"")
	initialPropValue = strings.TrimSuffix(initialPropValue, "'")

	content = strings.ReplaceAll(content, varNode.Name, initialPropValue)

	pageModel.Html.Content = content
}
