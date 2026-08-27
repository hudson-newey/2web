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
	"strings"

	"github.com/hudson-newey/2web/_shared/lists"
)

func NewReactiveEventNode(lexNodes []*lexer.V2LexNode) *reactiveEventNode {
	propName, err := scanners.NthToken(lexNodes, lexeme.TextContent, 1)
	if err != nil {
		panic(err)
	}

	reducer, err := scanners.NthToken(lexNodes, lexeme.TextContent, 2)
	if err != nil {
		panic(err)
	}

	assignmentSplit := strings.Split(reducer.Content, "=")
	assignmentSink := strings.TrimSpace(assignmentSplit[0])
	assignmentExpr := strings.TrimSpace(assignmentSplit[1])

	return &reactiveEventNode{
		eventName:      propName.Content,
		reducer:        reducer.Content,
		assignmentSink: assignmentSink,
		assignmentExpr: assignmentExpr,
	}
}

type reactiveEventNode struct {
	eventName string
	// A string of the entire reducer
	// $x = $y + $z + 1
	reducer string
	// Everything BEFORE the first equals sign (whitespace stripped)
	// $x
	assignmentSink string
	// Everything AFTER the first equals sign (whitespace trimmed)
	assignmentExpr string
	children       AbstractSyntaxTree
}

func (m *reactiveEventNode) Type() string {
	return "reactiveEventNode"
}

func (m *reactiveEventNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *reactiveEventNode) Content(page *page.Page, ast AbstractSyntaxTree) NodeContent {
	return NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

func (m *reactiveEventNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *reactiveEventNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

// What reactive variable this event binding is assigned to update
// In the following example, $x is the sink @click="$x = $y + $z + 1"
//
// TODO: We should add support for multiple sinks so that you can do something
// like @click="$x=$x+1; $y=$y*2;"
func (m *reactiveEventNode) reactiveVariableSink(ast AbstractSyntaxTree) *reactiveVariableNode {
	matched, err := lists.Find(ast.reactiveVariables(), func(x *reactiveVariableNode) bool {
		return x.selector() == m.assignmentSink
	})
	if err != nil {
		panic(fmt.Sprintf("could not find compiler variable '%s' for event %s", m.assignmentSink, m.eventName))
	}
	return matched
}

// What reactive variables are used in this event handler
// In the following example, $y and $z are the deps
// @click="$x = $y + $z + 1"
func (m *reactiveEventNode) reactiveVariableDeps(ast AbstractSyntaxTree) []*reactiveVariableNode {
	return lists.Filter(ast.reactiveVariables(), func(x *reactiveVariableNode) bool {
		return strings.Contains(m.assignmentExpr, x.selector())
	})
}
