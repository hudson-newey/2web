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

func NewScriptNode(lexNodes []*lexer.V2LexNode) *scriptNode {
	sourceNode, err := scanners.FirstToken(lexNodes, lexerTokens.ScriptSource)
	if err != nil {
		panic(err)
	}

	return &scriptNode{
		lexerNodes: lexNodes,
		content:    sourceNode.Content,
	}
}

type scriptNode struct {
	lexerNodes []*lexer.V2LexNode
	content    string
	children   ast.AbstractSyntaxTree
}

func (m *scriptNode) Type() string {
	return "scriptNode"
}

func (m *scriptNode) Children() ast.AbstractSyntaxTree {
	return m.children
}

func (m *scriptNode) AddChild(child ast.Node) {
	m.children = append(m.children, child)
}

func (m *scriptNode) RemoveChild(child ast.Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

func (m *scriptNode) Content(page *page.Page) ast.NodeContent {
	return ast.NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.FromContent(m.content),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}
