package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewStyleNode(lexNodes []*lexer.V2LexNode) *styleNode {
	sourceNode, err := scanners.FirstToken(lexNodes, lexerTokens.StyleSource)
	if err != nil {
		panic(err)
	}

	return &styleNode{
		lexerNodes: lexNodes,
		content:    sourceNode.Content,
	}
}

type styleNode struct {
	lexerNodes []*lexer.V2LexNode
	content    string
	children   ast.AbstractSyntaxTree
}

func (m *styleNode) Type() string {
	return "StyleNode"
}

func (m *styleNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (m *styleNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (m *styleNode) CssContent() *css.CSSFile {
	return css.FromContent(m.content)
}

func (m *styleNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.NewTwoScriptFile()
}

func (m *styleNode) Children() ast.AbstractSyntaxTree {
	return m.children
}

func (m *styleNode) AddChild(child ast.Node) {
	m.children = append(m.children, child)
}

func (m *styleNode) RemoveChild(child ast.Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}
