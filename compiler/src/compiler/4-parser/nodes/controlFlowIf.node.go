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

func NewControlFlowIfNode(lexNodes []*lexer.V2LexNode) *controlFlowIfNode {
	expression, err := scanners.NthToken(lexNodes, lexeme.TextContent, 1)
	if err != nil {
		panic(err)
	}

	content, err := scanners.NthToken(lexNodes, lexeme.TextContent, 2)
	if err != nil {
		panic(err)
	}

	return &controlFlowIfNode{
		expression: expression.Content,
		content:    content.Content,
		reactiveProp: &reactivePropertyNode{
			propName:      "hidden",
			reducer:       expression.Content,
			markupContent: ifExprNodeMarkup(expression.Content, content.Content),
		},
	}
}

type controlFlowIfNode struct {
	expression   string
	content      string
	children     AbstractSyntaxTree
	reactiveProp *reactivePropertyNode
}

func (m *controlFlowIfNode) Type() string {
	return "controlFlowIfNode"
}

func (m *controlFlowIfNode) Children() AbstractSyntaxTree {
	return AbstractSyntaxTree{m.reactiveProp}
}

func (m *controlFlowIfNode) MarkupContent() string {
	return m.reactiveProp.MarkupContent()
}

func (m *controlFlowIfNode) Content(page *page.Page, _ast AbstractSyntaxTree) NodeContent {
	return NodeContent{
		HtmlContent:      page.Html,
		TwoScriptContent: twoscript.NewTwoScriptFile(),
		CssContent:       css.NewCssFile(),
		JsContent:        javascript.NewJsFile(),
	}
}

func (m *controlFlowIfNode) AddChild(child Node) {
	m.reactiveProp.AddChild(child)
}

func (m *controlFlowIfNode) RemoveChild(child Node) {
	m.reactiveProp.RemoveChild(child)
}

func ifExprNodeMarkup(expression string, content string) string {
	// toggle the hidden global attribute to show/hide elements.
	// Kind of neat that this exists imo.
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/hidden
	return fmt.Sprintf(`<span *hidden="%s">%s</span>`, expression, content)
}
