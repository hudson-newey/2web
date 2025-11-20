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

// import importName from "importPath";
func NewScriptImportNode(lexNodes []*lexer.V2LexNode) *scriptImportNode {
	importNameNode, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 1)
	if err != nil {
		panic(err)
	}

	importPathNode, err := scanners.NthToken(lexNodes, lexerTokens.TextContent, 2)
	if err != nil {
		panic(err)
	}

	return &scriptImportNode{
		importName: importNameNode.Content,
		importPath: importPathNode.Content,
	}
}

type scriptImportNode struct {
	// What was the alias that the import was given?
	importName string

	// What is the path of the import?
	importPath string

	children ast.AbstractSyntaxTree
}

func (model *scriptImportNode) Type() string {
	return "importNode"
}

func (model *scriptImportNode) HtmlContent() *html.HTMLFile {
	return html.NewHtmlFile()
}

func (model *scriptImportNode) JsContent() *javascript.JSFile {
	return javascript.NewJsFile()
}

func (model *scriptImportNode) CssContent() *css.CSSFile {
	return css.NewCssFile()
}

func (model *scriptImportNode) TwoScriptContent() *twoscript.TwoScriptFile {
	return twoscript.NewTwoScriptFile()
}

func (model *scriptImportNode) Children() ast.AbstractSyntaxTree {
	return model.children
}

func (model *scriptImportNode) AddChild(child ast.Node) {
	model.children = append(model.children, child)
}

func (model *scriptImportNode) RemoveChild(child ast.Node) {
	for i, c := range model.children {
		if c == child {
			model.children = append(model.children[:i], model.children[i+1:]...)
			return
		}
	}
}
