package devtools

import (
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/models"
)

func InjectDevTools(pageContent string) string {
	devtoolsTemplateResult := createDevtoolsTemplate()
	return document.InjectContent(pageContent, devtoolsTemplateResult, document.Body)
}

func createDevtoolsTemplate() string {
	devtoolsHtml, err := document.BuildTemplate(devtoolsHtmlSource(), nil)
	if err != nil {
		// Handle the error, maybe add it to errorList
		documentErrors.AddError(models.Error{Message: "Failed to inject devtools template: " + err.Error()})
	}

	return devtoolsHtml
}
