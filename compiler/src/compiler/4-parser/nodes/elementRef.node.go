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

// While <code> nodes are normally regular HTML elements, I have decided to
// break this compatibility in 2web to make it easier to handle code content.
// In 2web, if you have a <code> block, everything inside will be automatically
// escaped into HTML escape sequences, meaning that if you include <style>,
// <script>, etc... tags inside of a <code> block, they will not be interpreted.
// By making this distinction at the node level (instead of using the tag name),
// it makes this distinction VERY clear without any ambiguity.
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
