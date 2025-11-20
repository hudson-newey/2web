package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewScriptReactiveVariableNode(lexNodes []*lexer.V2LexNode) *scriptReactiveVariableNode {
	variableName, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 1)
	if err != nil {
		panic(err)
	}

	initialValue, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 2)
	if err != nil {
		panic(err)
	}

	return &scriptReactiveVariableNode{
		variableName: variableName.Content,
		initialValue: initialValue.Content,
	}
}

type scriptReactiveVariableNode struct {
	variableName string
	initialValue string
	children     ast.AbstractSyntaxTree
}

func (model *scriptReactiveVariableNode) Type() string {
	return "reactiveVariableNode"
}

func (model *scriptReactiveVariableNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (model *scriptReactiveVariableNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (model *scriptReactiveVariableNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}

func (model *scriptReactiveVariableNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.NewTwoScriptFile()
}

func (model *scriptReactiveVariableNode) Children() ast.AbstractSyntaxTree {
	return model.children
}

func (model *scriptReactiveVariableNode) AddChild(child ast.Node) {
	model.children = append(model.children, child)
}

func (model *scriptReactiveVariableNode) RemoveChild(child ast.Node) {
	for i, c := range model.children {
		if c == child {
			model.children = append(model.children[:i], model.children[i+1:]...)
			return
		}
	}
}
