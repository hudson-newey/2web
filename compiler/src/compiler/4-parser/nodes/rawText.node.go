package nodes

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

// RawTextNode renders literal text without any interpretation.
//
// It is used to degrade gracefully on syntax errors: the raw source that the
// parser consumed is kept visible in the page (so the user can see what the
// compiler choked on) while the error overlay explains the problem.
type RawTextNode struct {
	content string
}

func NewRawTextNode(content string) *RawTextNode {
	return &RawTextNode{content: content}
}

func (m *RawTextNode) Type() string {
	return "RawTextNode"
}

func (m *RawTextNode) Children() AbstractSyntaxTree {
	return AbstractSyntaxTree{}
}

func (m *RawTextNode) MarkupContent() string {
	return html.EscapeHtml(m.content)
}

func (m *RawTextNode) Content(page *page.Page, _ *ReactiveIndex) NodeContent {
	return NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

func (m *RawTextNode) AddChild(child Node) {}

func (m *RawTextNode) RemoveChild(child Node) {}
