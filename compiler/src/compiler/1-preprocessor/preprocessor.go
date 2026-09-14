package preprocessor

import (
	"hudson-newey/2web/src/content/assets"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/content/xml"
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
		!txt.IsTxtFile(filePath) {

		// Components are expanded first so that the layout expansion and the
		// compilation pipeline can see the fully expanded page.
		//
		// The scope counter is per page (the component style scope ids are
		// deterministic per page), so it doesn't leak between the pages that
		// the build compiles in parallel.
		scopeCounter := 0
		ssgResult = expandComponents(filePath, ssgResult, map[string]bool{}, &scopeCounter)

		// Before we expand the HTML partials, we need to expand the layouts because
		// the layout may contain the doctype, html, head, and body tags that would
		// cause the partial expansion to fail.
		ssgResult = expandLayout(filePath, ssgResult)

		// The compile time virtual functions ($uid(), $env(), $readFile())
		// are expanded last, so that component content (which is inlined by
		// the component expansion) is covered too.
		ssgResult = expandVirtualFunctions(filePath, ssgResult)

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
