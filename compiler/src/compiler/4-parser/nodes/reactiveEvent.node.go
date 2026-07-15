package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewReactiveEventNode(lexNodes []*lexer.V2LexNode) *reactiveEventNode {
	propName, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 1)
	if err != nil {
		panic(err)
	}

	reducer, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 4)
	if err != nil {
		panic(err)
	}

	return &reactiveEventNode{
		eventName: propName.Content,
		reducer:   reducer.Content,
	}
}

type reactiveEventNode struct {
	eventName string
	reducer   string
	children  ast.AbstractSyntaxTree
}

func (m *reactiveEventNode) Type() string {
	return "reactiveEventNode"
}

func (m *reactiveEventNode) Children() ast.AbstractSyntaxTree {
	return m.children
}

func (m *reactiveEventNode) AddChild(child ast.Node) {
	m.children = append(m.children, child)
}

func (m *reactiveEventNode) RemoveChild(child ast.Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

func (m *reactiveEventNode) Content(page *page.Page) ast.NodeContent {
	return ast.NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}
