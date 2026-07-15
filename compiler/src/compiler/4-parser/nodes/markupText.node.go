package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewMarkupTextNode(lexNode []*lexer.V2LexNode) *markupTextNode {
	return &markupTextNode{
		content: lexNode[0].Content,
	}
}

type markupTextNode struct {
	content  string
	children ast.AbstractSyntaxTree
}

func (m *markupTextNode) Type() string {
	return "MarkupTextNode"
}

func (m *markupTextNode) Children() ast.AbstractSyntaxTree {
	return m.children
}

func (m *markupTextNode) AddChild(child ast.Node) {
	m.children = append(m.children, child)
}

func (m *markupTextNode) RemoveChild(child ast.Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

func (m *markupTextNode) Content(page *page.Page) ast.NodeContent {
	HtmlContent := html.FromContent(page.Html.Content + m.content)
	return ast.NodeContent{
		HtmlContent:      HtmlContent,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}
