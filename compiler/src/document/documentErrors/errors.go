package documentErrors

import (
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/models"
)

var errorList []models.Error

func AddError(errorModel models.Error) {
	errorList = append(errorList, errorModel)
}

func InjectErrors(pageContent string) string {
	if len(errorList) == 0 {
		return pageContent
	}

	errorTemplateResult := createErrorTemplate(errorList)

	return document.InjectContent(pageContent, errorTemplateResult, document.Body)
}

// creates a HTML error template that can be used to display errors
// in the browser
func createErrorTemplate(errors []models.Error) string {
	errorHtml, err := document.BuildTemplate(errorHtmlSource(), errors)
	if err != nil {
		// Handle the error, maybe add it to errorList
		AddError(models.Error{Message: "Failed to render error template: " + err.Error()})
	}

	return errorHtml
}
