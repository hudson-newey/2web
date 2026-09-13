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

func NewReactiveEventNode(lexNodes []*lexer.V2LexNode, context *ParseContext) Node {
	propName, err := scanners.NthToken(lexNodes, lexeme.TextContent, 1)
	if err != nil {
		return context.DegradedNode(
			"reactive event binding is missing an event name. Events are bound with '@event=\"...\"'",
			lexNodes,
		)
	}

	reducer, err := scanners.NthToken(lexNodes, lexeme.TextContent, 2)
	if err != nil {
		return context.DegradedNode(
			"reactive event binding is missing a reducer. Events are bound with '@event=\"reducer\"'",
			lexNodes,
		)
	}

	assignmentSink, assignmentExpr := parseReducer(reducer.Content)
	if assignmentSink == "" {
		return context.DegradedNode(
			"reactive event reducer must assign to a reactive variable (e.g. '$count = $count + 1' or '$count++')",
			lexNodes,
		)
	}

	markupContent := fmt.Sprintf(
		"@%s=\"%s",
		propName.Content,
		reducer.Content,
	)

	return &reactiveEventNode{
		eventName:      strings.TrimSpace(propName.Content),
		reducer:        strings.TrimSpace(reducer.Content),
		assignmentSink: assignmentSink,
		assignmentExpr: assignmentExpr,
		markupContent:  markupContent,
	}
}

// parseReducer splits an event reducer into its assignment sink and
// assignment expression.
//
// e.g. "$x = $y + 1" sinks to $x with the expression "$y + 1".
//
// The increment and decrement shorthand is expanded into an assignment:
// e.g. "$x++" sinks to $x with the expression "$x + 1".
func parseReducer(reducer string) (sink string, expression string) {
	assignmentSplit := strings.Split(reducer, "=")

	if len(assignmentSplit) >= 2 {
		return strings.TrimSpace(assignmentSplit[0]), strings.TrimSpace(assignmentSplit[1])
	}

	// Increment/decrement shorthand: "$x++" and "$x--".
	trimmed := strings.TrimSpace(assignmentSplit[0])

	if strings.HasSuffix(trimmed, "++") {
		sink = strings.TrimSpace(strings.TrimSuffix(trimmed, "++"))
		return sink, sink + " + 1"
	}

	if strings.HasSuffix(trimmed, "--") {
		sink = strings.TrimSpace(strings.TrimSuffix(trimmed, "--"))
		return sink, sink + " - 1"
	}

	return "", ""
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
	markupContent  string
	children       AbstractSyntaxTree
}

func (m *reactiveEventNode) Type() string {
	return "reactiveEventNode"
}

func (m *reactiveEventNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *reactiveEventNode) MarkupContent() string {
	return m.markupContent
}

func (m *reactiveEventNode) Content(page *page.Page, _ *ReactiveIndex) NodeContent {
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

func (m *reactiveEventNode) selector() string {
	// TODO: Use definitions from lexer here instead
	return fmt.Sprintf("@%s=\"%s\"", m.eventName, m.reducer)
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
