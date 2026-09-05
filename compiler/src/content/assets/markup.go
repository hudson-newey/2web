package assets

import (
	twoWeb "hudson-newey/2web/src/content/2web"
	"hudson-newey/2web/src/content/docx"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/odt"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/content/xml"
)

func IsMarkupFile(filePath string) bool {
	return html.IsHtmlFile(filePath) ||
		twoWeb.IsTwoWebFile(filePath) ||
		markdown.IsMarkdownFile(filePath) ||

		docx.IsDocxFile(filePath) ||
		odt.IsOdtFile(filePath) ||

		ShouldUseAssetPassthrough(filePath) ||

		// Most browsers can render xml natively, therefore, we treat them as markup
		// files.
		xml.IsXmlFile(filePath) ||

		// Most browsers treat .txt files as plain text, meaning that they can be
		// rendered natively.
		// Therefore, I treat them as markup files.
		txt.IsTxtFile(filePath)
}
