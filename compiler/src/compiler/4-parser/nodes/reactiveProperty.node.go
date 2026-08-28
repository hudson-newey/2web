package nodes

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
	"slices"
	"strings"

	"github.com/hudson-newey/2web/_shared/lists"
)

func NewReactivePropertyNode(lexNodes []*lexer.V2LexNode) *reactivePropertyNode {
	propName, err := scanners.NthToken(lexNodes, lexeme.TextContent, 1)
	if err != nil {
		panic(err)
	}

	reducer, err := scanners.NthToken(lexNodes, lexeme.TextContent, 4)
	if err != nil {
		panic(err)
	}

	markupContent := fmt.Sprintf(
		"*%s=\"%s\"",
		propName.Content,
		reducer.Content,
	)

	return &reactivePropertyNode{
		propName:      propName.Content,
		reducer:       reducer.Content,
		markupContent: markupContent,
	}
}

type reactivePropertyNode struct {
	propName      string
	reducer       string
	markupContent string
	children      AbstractSyntaxTree
}

func (m *reactivePropertyNode) Type() string {
	return "reactivePropertyNode"
}

func (m *reactivePropertyNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *reactivePropertyNode) MarkupContent() string {
	return m.markupContent
}

func (m *reactivePropertyNode) Content(page *page.Page, _ast AbstractSyntaxTree) NodeContent {
	return NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

func (m *reactivePropertyNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *reactivePropertyNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

func (m *reactivePropertyNode) selector() string {
	// TODO: Use definitions from lexer here instead
	return fmt.Sprintf("*%s=\"%s\"", m.propName, m.reducer)
}

// Some properties such as text assignment can be evaluated at compile time
// meaning that we don't need to ship any runtime code.
func (m *reactivePropertyNode) canCompilerInline() bool {
	supportedPropSinks := []string{"textContent", "innerText"}
	return slices.Contains(supportedPropSinks, m.propName)
}

// Finds all reactive variable dependencies for a property binding
func (m *reactivePropertyNode) reactiveVariableDeps(ast AbstractSyntaxTree) []*reactiveVariableNode {
	return lists.Filter(ast.reactiveVariables(), func(x *reactiveVariableNode) bool {
		return strings.Contains(m.reducer, x.selector())
	})
}
