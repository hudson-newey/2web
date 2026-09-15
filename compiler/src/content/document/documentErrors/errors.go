package documentErrors

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/models"
	"sync"
	"time"
)

type errorTemplateData struct {
	Errors    []*models.Error
	CreatedAt string
}

// totalErrors is appended to by all concurrently compiled pages, so appends
// must be synchronized to prevent lost errors and slice header corruption.
var totalErrorsMutex sync.Mutex
var totalErrors []*models.Error

func AddErrors(errorModels ...*models.Error) {
	totalErrorsMutex.Lock()
	defer totalErrorsMutex.Unlock()

	totalErrors = append(totalErrors, errorModels...)
}

func IsErrorFree() bool {
	totalErrorsMutex.Lock()
	defer totalErrorsMutex.Unlock()

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

	for _, errorModel := range snapshotErrors() {
		errorModel.PrintError()
	}
}

func snapshotErrors() []*models.Error {
	totalErrorsMutex.Lock()
	defer totalErrorsMutex.Unlock()

	return totalErrors
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
		// The error overlay failing to render must not kill the build; the
		// error is still reported in the terminal.
		return "<!-- 2web: failed to render the compiler error overlay -->"
	}

	return errorHtml
}
