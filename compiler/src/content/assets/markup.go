package assets

import (
	twoWeb "hudson-newey/2web/src/content/2web"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/xhtml"
	"hudson-newey/2web/src/content/xml"
	"hudson-newey/2web/src/content/xslt"
)

func IsMarkupFile(filePath string) bool {
	return html.IsHtmlFile(filePath) ||
		xhtml.IsXhtmlFile(filePath) ||
		twoWeb.IsTwoWebFile(filePath) ||
		markdown.IsMarkdownFile(filePath) ||

		// Most browsers can render xml natively, therefore, we treat them as markup
		// files.
		xml.IsXmlFile(filePath) ||
		xslt.IsXsltFile(filePath)
}
