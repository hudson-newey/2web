package imports

import (
	"fmt"
	"hudson-newey/2web/src/content/document/documentErrors"
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

func expandImport(
	workingPath string,
	content string,
	component *models.Component,
) string {
	componentTemplate, err := component.HtmlContent(workingPath)
	if err != nil {
		documentErrors.AddError(models.Error{
			FilePath: workingPath,
			Message:  err.Error(),
		})

		// We return the input content without modification in the hopes that the
		// page will be semi-functional and assist in debugging how the error
		// occurred.
		return content
	}

	templateSelector := fmt.Sprintf("<%s />", component.DomSelector)
	result := strings.ReplaceAll(content, templateSelector, componentTemplate)

	return result
}
