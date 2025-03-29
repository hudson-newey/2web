package documentErrors

import (
	"bytes"
	"html/template"
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

	return pageContent + createErrorTemplate(errorList)
}

// creates a HTML error template that can be used to display errors
// in the browser
func createErrorTemplate(errors []models.Error) string {
	htmlTemplate := errorHtmlTemplate()

	errorTemplate := template.Must(template.New("errorTemplate").Parse(htmlTemplate))
	var buf bytes.Buffer
	err := errorTemplate.Execute(&buf, errors)
	if err != nil {
		// Handle the error, maybe add it to errorList
		AddError(models.Error{Message: "Failed to render error template: " + err.Error()})
	}
	errorHtml := buf.String()

	return errorHtml
}
