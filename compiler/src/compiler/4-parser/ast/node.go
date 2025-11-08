package ast

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

type nodeType = string

type Node interface {
	Type() nodeType

	HtmlContent() *html.HTMLFile
	JsContent() *javascript.JSFile
	CssContent() *css.CSSFile
}

func HasHtmlContent(node Node) bool {
	return node.HtmlContent().Content != ""
}

func HasJsContent(node Node) bool {
	return node.JsContent().Content != ""
}

func HasCssContent(node Node) bool {
	return node.CssContent().Content != ""
}
