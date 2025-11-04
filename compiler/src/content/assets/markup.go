package assets

import (
	twoWeb "hudson-newey/2web/src/content/2web"
	"hudson-newey/2web/src/content/docx"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/odt"
	"hudson-newey/2web/src/content/pdf"
	"hudson-newey/2web/src/content/txt"
	"hudson-newey/2web/src/content/xhtml"
	"hudson-newey/2web/src/content/xml"
	"hudson-newey/2web/src/content/xslt"
)

func IsMarkupFile(filePath string) bool {
	return html.IsHtmlFile(filePath) ||
		xhtml.IsXhtmlFile(filePath) ||
		twoWeb.IsTwoWebFile(filePath) ||
		markdown.IsMarkdownFile(filePath) ||

		docx.IsDocxFile(filePath) ||
		odt.IsOdtFile(filePath) ||

		// PDF files are "technically" binary files, however, they act more like
		// markup because they can be rendered natively in most browsers.
		// Therefore, we treat them as markup files instead of including them as a
		// static asset because static assets have some special behavior like name
		// mangling and lazy loading that do not make sense for PDFs.
		pdf.IsPdfFile(filePath) ||

		// Most browsers can render xml natively, therefore, we treat them as markup
		// files.
		xml.IsXmlFile(filePath) ||
		xslt.IsXsltFile(filePath) ||

		// Most browsers treat .txt files as plain text, meaning that they can be
		// rendered natively.
		// Therefore, I treat them as markup files.
		txt.IsTxtFile(filePath)
}
