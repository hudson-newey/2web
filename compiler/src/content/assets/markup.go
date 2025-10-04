package assets

import (
	twoWeb "hudson-newey/2web/src/content/2web"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/xhtml"
)

func IsMarkupFile(filePath string) bool {
	return html.IsHtmlFile(filePath) ||
		xhtml.IsXhtmlFile(filePath) ||
		twoWeb.IsTwoWebFile(filePath) ||
		markdown.IsMarkdownFile(filePath)
}
