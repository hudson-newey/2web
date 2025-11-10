package preprocessor

import (
	"hudson-newey/2web/src/content/assets"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/content/xhtml"
	"hudson-newey/2web/src/content/xml"
	"hudson-newey/2web/src/content/xslt"
)

func ProcessStaticSite(filePath string, content string, expandPartials bool) string {
	ssgResult := content

	// We convert markdown files into HTML first so that layouts and partials can
	// be applied to the resulting HTML.
	if markdown.IsMarkdownFile(filePath) {
		markdownFile := markdown.MarkdownFile{
			Content: ssgResult,
		}
		ssgResult = markdownFile.ToHtml().Content
	}

	if assets.IsMarkupFile(filePath) &&
		!xml.IsXmlFile(filePath) &&
		!xslt.IsXsltFile(filePath) &&
		!txt.IsTxtFile(filePath) {
		// Before we expand the HTML partials, we need to expand the layouts because
		// the layout may contain the doctype, html, head, and body tags that would
		// cause the partial expansion to fail.
		//
		// TODO: Add support for markdown & xhtml layouts.
		if !xhtml.IsXhtmlFile(filePath) && !markdown.IsMarkdownFile(filePath) {
			buffered := []byte(ssgResult)
			expandLayout(filePath, buffered)
		}

		// 2Web supports partial content, meaning that pages don't need and doctype,
		// html, head, meta, or body tags.
		// The user can just start writing the pages content, and the compiler can
		// figure out what should be in the body vs head.
		if expandPartials {
			ssgResult = html.ExpandPartial(ssgResult)
		}
	}

	return ssgResult
}
