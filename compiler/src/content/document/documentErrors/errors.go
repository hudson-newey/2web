package documentErrors

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/models"
	"time"
)

type errorTemplateData struct {
	Errors    []*models.Error
	CreatedAt string
}

var totalErrors []*models.Error

func AddErrors(errorModels ...*models.Error) {
	totalErrors = append(totalErrors, errorModels...)
}

func IsErrorFree() bool {
	return len(totalErrors) == 0
}

// TODO: Use a page reference here
func InjectErrors(content string, pageErrors []*models.Error) string {
	errorTemplateResult := createErrorTemplate(pageErrors)
	return document.InjectContent(content, errorTemplateResult, document.BodyTop)
}

func PrintDocumentErrors() {
	if cli.GetArgs().IsSilent {
		return
	}

	for _, errorModel := range totalErrors {
		errorModel.PrintError()
	}
}

// creates a HTML error template that can be used to display errors
// in the browser
func createErrorTemplate(errors []*models.Error) string {
	creationTime := time.Now().Format(time.DateTime)

	templateData := errorTemplateData{
		Errors:    errors,
		CreatedAt: creationTime,
	}

	errorHtml, err := document.BuildTemplate(errorHtmlSource(), templateData)
	if err != nil {
		panic(err)
	}

	return errorHtml
}
