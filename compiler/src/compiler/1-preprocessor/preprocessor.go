package preprocessor

import (
	"hudson-newey/2web/src/content/assets"
	"hudson-newey/2web/src/content/docx"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/odt"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/content/xml"
	"hudson-newey/2web/src/content/xslt"
)

func ProcessStaticSite(filePath string, content string) string {
	ssgResult := content

	if assets.IsMarkupFile(filePath) &&
		!markdown.IsMarkdownFile(filePath) &&
		!xml.IsXmlFile(filePath) &&
		!xslt.IsXsltFile(filePath) &&
		!txt.IsTxtFile(filePath) {
		// 2Web supports partial content, meaning that pages don't need and doctype,
		// html, head, meta, or body tags.
		// The user can just start writing the pages content, and the compiler can
		// figure out what should be in the body vs head.
		ssgResult = html.ExpandPartial(content)
	} else if markdown.IsMarkdownFile(filePath) {
		markdownFile := markdown.MarkdownFile{
			Content: content,
		}
		ssgResult = markdownFile.ToHtml().Content

		// Markdown files are typically compiled as html partials. Developers
		// typically (and shouldn't) declare a doctype, header, etc...
		// therefore, we also expand html partials once the html document has been
		// created.
		ssgResult = html.ExpandPartial(ssgResult)
	} else if docx.IsDocxFile(filePath) {
		ssgResult = html.ExpandPartial(content)
	} else if odt.IsOdtFile(filePath) {
		ssgResult = html.ExpandPartial(content)
	}

	return ssgResult
}
