package nodes

import (
	"fmt"
	"strings"

	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewControlFlowIfNode(lexNodes []*lexer.V2LexNode) *controlFlowIfNode {
	// The grammar captures everything between the condition's parentheses and
	// the body's curly braces, so the expression and content can contain
	// whitespace and multiple tokens.
	expression, err := scanners.CapturedContent(lexNodes, lexeme.BracketOpen, lexeme.BracketClosed)
	if err != nil {
		panic(err)
	}

	content, err := scanners.CapturedContent(lexNodes, lexeme.CurlyOpen, lexeme.CurlyClosed)
	if err != nil {
		panic(err)
	}

	expression = strings.TrimSpace(expression)

	// The markup pass renders the reactive property node (which owns the
	// conditionally rendered markup) through Children(), so the if node itself
	// must not render any markup. Rendering both would duplicate the content.
	return &controlFlowIfNode{
		expression: expression,
		content:    content,
		reactiveProp: &reactivePropertyNode{
			propName: "hidden",
			// Keep the raw expression as the reducer so that reactive variable
			// dependency detection can find the variables used in the condition.
			reducer:       expression,
			markupContent: ifExprNodeMarkup(expression, content),
			// `hidden` hides an element when it is true, but `@if (condition)`
			// should render its content when the condition is true. Negate the
			// condition so that the content is shown when the condition is
			// truthy.
			negated: true,
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
	// The conditionally rendered markup is owned by the reactive property node
	// (see Children), which is rendered during the markup pass.
	return ""
}

func (m *controlFlowIfNode) Content(page *page.Page, _ *ReactiveIndex) NodeContent {
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
