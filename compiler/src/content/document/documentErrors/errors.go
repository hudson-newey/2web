package documentErrors

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/models"
)

var pageErrors []models.Error
var totalErrors []models.Error

func AddErrors(errorModels ...models.Error) {
	for _, err := range errorModels {
		pageErrors = append(pageErrors, err)
		totalErrors = append(totalErrors, err)
	}
}

// TODO: errors should be attached to the page model instead of here
func ResetPageErrors() {
	pageErrors = []models.Error{}
}

func IsErrorFree() bool {
	return len(totalErrors) == 0
}

func InjectErrors(pageContent string) string {
	if len(pageErrors) == 0 {
		return pageContent
	}

	errorTemplateResult := createErrorTemplate(pageErrors)

	return document.InjectContent(pageContent, errorTemplateResult, document.Body)
}

func PrintDocumentErrors() {
	for _, errorModel := range totalErrors {
		cli.PrintError(errorModel)
	}
}

// creates a HTML error template that can be used to display errors
// in the browser
func createErrorTemplate(errors []models.Error) string {
	errorHtml, err := document.BuildTemplate(errorHtmlSource(), errors)
	if err != nil {
		// Handle the error, maybe add it to errorList
		AddErrors(models.Error{Message: "failed to render error template: " + err.Error()})
	}

	return errorHtml
}
