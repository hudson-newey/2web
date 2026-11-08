package builder

import (
	"hudson-newey/2web/src/content"
	"hudson-newey/2web/src/content/assets"
	"hudson-newey/2web/src/content/page"
)

func buildFromBinary(
	inputPath string,
	data []byte,
	isFullPage bool,
) (page.Page, bool) {
	if assets.ShouldUseAssetPassthrough(inputPath) {
		return buildStaticAsset(inputPath, data)
	}

	compiledPage := page.NewPage()
	compiledPage.Html.AddContent("[Unsupported Binary Format]")

	return compiledPage, false
}

func buildStaticAsset(inputPath string, data []byte) (page.Page, bool) {
	pageModel := page.NewPage()
	var pdfModel content.BinaryFile = assets.FromContent(inputPath, data)

	pageModel.AddAsset(&pdfModel)

	return pageModel, true
}
