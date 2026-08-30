package nodes

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/page"
)

func NewTextOutputNode(lexNodes []*lexer.V2LexNode) *textOutputNode {
	expression, err := scanners.NthToken(lexNodes, lexeme.TextContent, 1)
	if err != nil {
		panic(err)
	}

	// so... text nodes are really just a short hand for a reactive property
	// node's textContent property...
	// This is a big weird in the compiler. And it's the only time that a
	// constructor doesn't return an instance of itself.
	// Although, I do think it's quite a neat solution to the problem since it
	// allows us to re-use the reactive property & compiler logic without
	// needing to put text output node expansion in the pre-processor in a
	// separate "mini compiler".
	return &textOutputNode{
		expression: expression.Content,
		reactiveProp: &reactivePropertyNode{
			propName:      "textContent",
			reducer:       expression.Content,
			markupContent: textNodeMarkup(expression.Content),
		},
	}
}

type textOutputNode struct {
	// e.g. in {{ $count * 2}} the expression would be $count * 2
	expression   string
	reactiveProp *reactivePropertyNode
}

func (m *textOutputNode) Type() string {
	return "textOutputNode"
}

func (m *textOutputNode) Children() AbstractSyntaxTree {
	return AbstractSyntaxTree{m.reactiveProp}
}

func (m *textOutputNode) MarkupContent() string {
	return m.reactiveProp.MarkupContent()
}

func (m *textOutputNode) Content(page *page.Page, ast AbstractSyntaxTree) NodeContent {
	// A kinda cool solution which lets us re-use the reactive prop content
	// which hooks up event listeners, etc...
	return m.reactiveProp.Content(page, ast)
}

func (m *textOutputNode) AddChild(child Node) {
	m.reactiveProp.AddChild(child)
}

func (m *textOutputNode) RemoveChild(child Node) {
	m.reactiveProp.RemoveChild(child)
}

func textNodeMarkup(expression string) string {
	return fmt.Sprintf(`<span *textContent="%s"></span>`, expression)
}
