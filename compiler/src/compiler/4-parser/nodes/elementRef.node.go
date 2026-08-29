package nodes

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewElementRefNode(lexNodes []*lexer.V2LexNode) *elementRefNode {
	id, err := scanners.NthToken(lexNodes, lexeme.TextContent, 1)
	if err != nil {
		panic(err)
	}

	return &elementRefNode{
		id: id.Content,
	}
}

type elementRefNode struct {
	id       string
	children AbstractSyntaxTree
}

func (m *elementRefNode) Type() string {
	return "elementRefNode"
}

func (m *elementRefNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *elementRefNode) MarkupContent() string {
	return fmt.Sprintf(`id="%s"`, m.id)
}

func (m *elementRefNode) Content(page *page.Page, _ast AbstractSyntaxTree) NodeContent {
	return NodeContent{
		HtmlContent:      page.Html,
		TwoScriptContent: twoscript.NewTwoScriptFile(),
		CssContent:       css.NewCssFile(),
		JsContent:        javascript.NewJsFile(),
	}
}

func (m *elementRefNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *elementRefNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}
