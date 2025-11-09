package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewStyleNode(lexNodes []*lexer.V2LexNode) *styleNode {
	// Find the lexNode that is a StyleSource token
	var content string
	for _, lexNode := range lexNodes {
		if lexNode.Token == lexerTokens.StyleSource {
			content = lexNode.Content
			break
		}
	}

	return &styleNode{
		lexerNodes: lexNodes,
		content:    content,
	}
}

type styleNode struct {
	lexerNodes []*lexer.V2LexNode
	content    string
}

func (model *styleNode) Type() string {
	return "StyleNode"
}

func (model *styleNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (model *styleNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (model *styleNode) CssContent() *css.CSSFile {
	return css.FromContent(model.content)
}

func (model *styleNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.NewTwoScriptFile()
}
