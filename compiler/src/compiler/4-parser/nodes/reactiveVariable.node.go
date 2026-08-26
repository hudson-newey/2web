package nodes

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewScriptReactiveVariableNode(lexNodes []*lexer.V2LexNode) *scriptReactiveVariableNode {
	variableName, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 1)
	if err != nil {
		panic(err)
	}

	initialValue, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 2)
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
	children     AbstractSyntaxTree
}

func (m *scriptReactiveVariableNode) Type() string {
	return "reactiveVariableNode"
}

func (m *scriptReactiveVariableNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *scriptReactiveVariableNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *scriptReactiveVariableNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

func (m *scriptReactiveVariableNode) Content(page *page.Page, _ast AbstractSyntaxTree) NodeContent {
	return NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}
