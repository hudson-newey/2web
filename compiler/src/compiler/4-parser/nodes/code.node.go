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
}

func (model *codeNode) Type() string {
	return "codeNode"
}

func (model *codeNode) HtmlContent() *html.HTMLFile {
	escapedContent := html.EscapeHtml(model.content)
	withCodeTags := model.startingCodeTagContent + escapedContent + "</code>"

	return html.FromContent(withCodeTags)
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
