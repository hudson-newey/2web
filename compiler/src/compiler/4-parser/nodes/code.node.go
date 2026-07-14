package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewCodeNode(lexNodes []*lexer.V2LexNode) *codeNode {
	startingCodeTagContent := ""

	// Find the lexNode that is a StyleSource token
	var content string
	for _, lexNode := range lexNodes {
		if lexNode.Token == lexerTokens.CodeSource {
			content = lexNode.Content
			break
		}

		// Find all content between the < and the > tokens (inclusive) so that
		// attributes on the code tag are preserved.
		startingCodeTagContent += lexNode.Content
	}

	return &codeNode{
		startingCodeTagContent: startingCodeTagContent,
		lexerNodes:             lexNodes,
		content:                content,
	}
}

// While <code> nodes are normally regular HTML elements, I have decided to
// break this compatibility in 2web to make it easier to handle code content.
// In 2web, if you have a <code> block, everything inside will be automatically
// escaped into HTML escape sequences, meaning that if you include <style>,
// <script>, etc... tags inside of a <code> block, they will not be interpreted.
// By making this distinction at the node level (instead of using the tag name),
// it makes this distinction VERY clear without any ambiguity.
type codeNode struct {
	startingCodeTagContent string
	lexerNodes             []*lexer.V2LexNode
	content                string
	children               ast.AbstractSyntaxTree
}

func (m *codeNode) Type() string {
	return "codeNode"
}

func (m *codeNode) escapedHtml() string {
	escapedContent := html.EscapeHtml(m.content)
	return m.startingCodeTagContent + escapedContent + "</code>"
}

func (m *codeNode) Children() ast.AbstractSyntaxTree {
	return m.children
}

func (m *codeNode) AddChild(child ast.Node) {
	m.children = append(m.children, child)
}

func (m *codeNode) RemoveChild(child ast.Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

func (m *codeNode) Content(page *page.Page) ast.NodeContent {
	HtmlContent := html.FromContent(page.Html.Content + m.escapedHtml())

	return ast.NodeContent{
		HtmlContent:      HtmlContent,
		TwoScriptContent: twoscript.NewTwoScriptFile(),
		CssContent:       css.NewCssFile(),
		JsContent:        javascript.NewJsFile(),
	}
}
