package builder

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/docx"
	"hudson-newey/2web/src/content/odt"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/pdf"
	"hudson-newey/2web/src/convert"
	"hudson-newey/2web/src/models"
)

func buildFromBinary(inputPath string, data []byte, isFullPage bool) (page.Page, bool) {
	if docx.IsDocxFile(inputPath) {
		return buildDocx(inputPath, isFullPage)
	} else if odt.IsOdtFile(inputPath) {
		return buildOdt(inputPath, isFullPage)
	} else if pdf.IsPdfFile(inputPath) {
		return buildPdf(inputPath)
	}

	compiledPage := page.NewPage()
	compiledPage.Html.AddContent("[Unsupported Binary Format]")

	return compiledPage, false
}

func buildDocx(inputPath string, isFullPage bool) (page.Page, bool) {
	docxModel := docx.NewDocxFile(inputPath)
	htmlContent, err := convert.ConvertFormat(docxModel.Data, "docx", "html")
	if err != nil {
		documentErrors.AddErrors(
			models.NewError(
				"Failed to convert 'docx' file to 'html'",
				inputPath,
				lexer.Position{},
			),
		)

		return page.NewPage(), false
	}

	// We use buildFromString to process the resulting HTML content
	// so that the reactive compiler and other features such as the devtools,
	// runtime optimizer, element refs, etc... still work.
	page, isErrorFree := buildFromString(inputPath, string(htmlContent), isFullPage)

	return page, isErrorFree
}

func buildOdt(inputPath string, isFullPage bool) (page.Page, bool) {
	odtModel := odt.NewOdtFile(inputPath)
	htmlContent, err := convert.ConvertFormat(odtModel.Data, "odt", "html")
	if err != nil {
		documentErrors.AddErrors(
			models.NewError(
				"Failed to convert 'odt' file to 'html'",
				inputPath,
				lexer.Position{},
			),
		)

		return page.NewPage(), false
	}

	page, isErrorFree := buildFromString(inputPath, string(htmlContent), isFullPage)

	return page, isErrorFree
}

func buildPdf(inputPath string) (page.Page, bool) {
	pageModel := page.NewPage()
	var pdfModel content.BinaryFile = pdf.NewPdfFile(inputPath)

	pageModel.AddAsset(&pdfModel)

	return pageModel, true
}
