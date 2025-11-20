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
	content  string
	children ast.AbstractSyntaxTree
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
	return model.children
}

func (model *markupTextNode) AddChild(child ast.Node) {
	model.children = append(model.children, child)
}

func (model *markupTextNode) RemoveChild(child ast.Node) {
	for i, c := range model.children {
		if c == child {
			model.children = append(model.children[:i], model.children[i+1:]...)
			return
		}
	}
}
