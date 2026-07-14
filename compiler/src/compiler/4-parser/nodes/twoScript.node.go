package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewTwoScriptNode(lexNodes []*lexer.V2LexNode) *twoScriptNode {
	// Find the lexNode that is a StyleSource token
	var content string
	for _, lexNode := range lexNodes {
		if lexNode.Token == lexerTokens.CompiledScriptSource {
			content = lexNode.Content
			break
		}
	}

	return &twoScriptNode{
		lexerNodes: lexNodes,
		content:    content,
	}
}

type twoScriptNode struct {
	lexerNodes []*lexer.V2LexNode
	content    string
	children   ast.AbstractSyntaxTree
}

func (m *twoScriptNode) Type() string {
	return "twoScriptNode"
}

func (m *twoScriptNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (m *twoScriptNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (m *twoScriptNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}

func (m *twoScriptNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.FromContent(m.content)
}

func (m *twoScriptNode) Children() ast.AbstractSyntaxTree {
	return m.children
}

func (m *twoScriptNode) AddChild(child ast.Node) {
	m.children = append(m.children, child)
}

func (m *twoScriptNode) RemoveChild(child ast.Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}
