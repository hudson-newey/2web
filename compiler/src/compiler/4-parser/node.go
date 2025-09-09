package parser

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

type Node interface {
	Tokens() []string

	HtmlContent() *html.HTMLFile
	JsContent() *javascript.JSFile
	CssContent() *css.CSSFile
}
