package parser

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

type MarkupTextNode struct {
	content string
}

func (model *MarkupTextNode) Tokens() []string {
	return []string{model.content}
}

func (model *MarkupTextNode) HtmlContent() *html.HTMLFile {
	return &html.HTMLFile{Content: model.content}
}

func (model *MarkupTextNode) JsContent() *javascript.JSFile {
	return &javascript.JSFile{Content: ""}
}

func (model *MarkupTextNode) CssContent() *css.CSSFile {
	return &css.CSSFile{Content: ""}
}
