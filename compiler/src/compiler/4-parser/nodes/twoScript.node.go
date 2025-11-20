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

func (model *twoScriptNode) Type() string {
	return "twoScriptNode"
}

func (model *twoScriptNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (model *twoScriptNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (model *twoScriptNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}

func (model *twoScriptNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.FromContent(model.content)
}

func (model *twoScriptNode) Children() ast.AbstractSyntaxTree {
	return model.children
}

func (model *twoScriptNode) AddChild(child ast.Node) {
	model.children = append(model.children, child)
}

func (model *twoScriptNode) RemoveChild(child ast.Node) {
	for i, c := range model.children {
		if c == child {
			model.children = append(model.children[:i], model.children[i+1:]...)
			return
		}
	}
}
