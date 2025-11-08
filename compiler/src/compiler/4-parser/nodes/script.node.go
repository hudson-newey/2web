package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

func NewScriptNode(lexNodes []*lexer.V2LexNode) *scriptNode {
	// Find the lexNode that is a StyleSource token
	var content string
	for _, lexNode := range lexNodes {
		if lexNode.Token == lexerTokens.ScriptSource {
			content = lexNode.Content
			break
		}
	}

	return &scriptNode{
		lexerNodes: lexNodes,
		content:    content,
	}
}

type scriptNode struct {
	lexerNodes []*lexer.V2LexNode
	content    string
}

func (model *scriptNode) Type() string {
	return "scriptNode"
}

func (model *scriptNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (model *scriptNode) JsContent() *javascript.JSFile {
	return javascript.FromContent(model.content)
}

func (model *scriptNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}
