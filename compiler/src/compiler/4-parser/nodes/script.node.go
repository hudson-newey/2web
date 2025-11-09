package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
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

func (model *scriptNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.NewTwoScriptFile()
}
