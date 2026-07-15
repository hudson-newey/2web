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

func NewReactivePropertyNode(lexNodes []*lexer.V2LexNode) *reactivePropertyNode {
	propName, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 1)
	if err != nil {
		panic(err)
	}

	reducer, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 2)
	if err != nil {
		panic(err)
	}

	return &reactivePropertyNode{
		propName: propName.Content,
		reducer:  reducer.Content,
	}
}

type reactivePropertyNode struct {
	propName string
	reducer  string
	children ast.AbstractSyntaxTree
}

func (m *reactivePropertyNode) Type() string {
	return "reactivePropertyNode"
}

func (m *reactivePropertyNode) Children() ast.AbstractSyntaxTree {
	return m.children
}

func (m *reactivePropertyNode) AddChild(child ast.Node) {
	m.children = append(m.children, child)
}

func (m *reactivePropertyNode) RemoveChild(child ast.Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

func (m *reactivePropertyNode) Content(page *page.Page) ast.NodeContent {
	return ast.NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}
