package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewMarkupTextNode(lexNode []*lexer.V2LexNode) *markupTextNode {
	return &markupTextNode{
		content: lexNode[0].Content,
	}
}

type markupTextNode struct {
	content string
}

func (model *markupTextNode) Type() string {
	return "MarkupTextNode"
}

func (model *markupTextNode) HtmlContent() *html.HTMLFile {
	return html.FromContent(model.content)
}

func (model *markupTextNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (model *markupTextNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}

func (model *markupTextNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.NewTwoScriptFile()
}

func (model *markupTextNode) Children() ast.AbstractSyntaxTree {
	return ast.AbstractSyntaxTree{}
}
