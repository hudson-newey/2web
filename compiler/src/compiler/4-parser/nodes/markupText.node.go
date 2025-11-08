package nodes

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

func NewMarkupTextNode(content string) *markupTextNode {
	return &markupTextNode{
		content: content,
	}
}

type markupTextNode struct {
	content string
}

func (model *markupTextNode) Tokens() []string {
	return []string{model.content}
}

func (model *markupTextNode) HtmlContent() *html.HTMLFile {
	return &html.HTMLFile{Content: model.content}
}

func (model *markupTextNode) JsContent() *javascript.JSFile {
	return &javascript.JSFile{Content: ""}
}

func (model *markupTextNode) CssContent() *css.CSSFile {
	return &css.CSSFile{Content: ""}
}
