package builder

import (
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/docx"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/convert"
	"hudson-newey/2web/src/models"
)

func buildFromBinary(inputPath string, data []byte) (page.Page, bool) {
	if docx.IsDocxFile(inputPath) {
		return buildDocx(inputPath)
	}

	compiledPage := page.NewPage()
	compiledPage.Html.AddContent("[Unsupported Binary Format]")

	return compiledPage, false
}

func buildDocx(inputPath string) (page.Page, bool) {
	docxModel := docx.NewDocxFile(inputPath)
	htmlContent, err := convert.ConvertFormat(docxModel.Data, "docx", "html")
	if err != nil {
		documentErrors.AddErrors(
			models.Error{
				FilePath: inputPath,
				Message:  "Failed to convert 'docx' file to 'html'",
			},
		)

		return page.NewPage(), false
	}

	// We use buildFromString to process the resulting HTML content
	// so that the reactive compiler and other features such as the devtools,
	// runtime optimizer, element refs, etc... still work.
	page, isErrorFree := buildFromString(inputPath, string(htmlContent))

	return page, isErrorFree
}
