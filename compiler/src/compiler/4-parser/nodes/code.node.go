package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewCodeNode(lexNodes []*lexer.V2LexNode) *codeNode {
	// Find the lexNode that is a StyleSource token
	var content string
	for _, lexNode := range lexNodes {
		if lexNode.Token == lexerTokens.ScriptSource {
			content = lexNode.Content
			break
		}
	}

	return &codeNode{
		lexerNodes: lexNodes,
		content:    content,
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
	lexerNodes []*lexer.V2LexNode
	content    string
}

func (model *codeNode) Type() string {
	return "codeNode"
}

func (model *codeNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (model *codeNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (model *codeNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}

func (model *codeNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.NewTwoScriptFile()
}
