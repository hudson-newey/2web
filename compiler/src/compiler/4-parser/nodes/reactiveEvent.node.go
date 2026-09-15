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
	if assignmentSink == "" && assignmentExpr == "" {
		return context.DegradedNode(
			"reactive event reducer must assign to a reactive variable or call an imported server function (e.g. '$count = $count + 1' or 'save($name)')",
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
		// The position of the binding is kept so that errors about the reducer
		// (e.g. a call to a function that wasn't imported from a server
		// script) can be reported against the binding.
		position: positionOf(lexNodes),
	}
}

// parseReducer splits an event reducer into its assignment sink and
// assignment expression.
//
// e.g. "$x = $y + 1" sinks to $x with the expression "$y + 1".
//
// The increment and decrement shorthand is expanded into an assignment:
// e.g. "$x++" sinks to $x with the expression "$x + 1".
//
// A reducer that is a direct function call (e.g. "save($name)") doesn't
// assign to anything. It is returned with an empty sink and the call as its
// expression, and is compiled into a standalone event listener that calls the
// function (see CompileServerCalls for the server function rpc calls).
func parseReducer(reducer string) (sink string, expression string) {
	trimmed := strings.TrimSpace(reducer)

	// Increment/decrement shorthand: "$x++" and "$x--".
	if strings.HasSuffix(trimmed, "++") {
		sink = strings.TrimSpace(strings.TrimSuffix(trimmed, "++"))
		return sink, sink + " + 1"
	}

	if strings.HasSuffix(trimmed, "--") {
		sink = strings.TrimSpace(strings.TrimSuffix(trimmed, "--"))
		return sink, sink + " - 1"
	}

	// A trailing semicolon is optional.
	trimmed = strings.TrimSuffix(trimmed, ";")

	// An assignment assigns to the reactive variable before the first "=".
	// The split must happen at the FIRST "=" only (and must not treat
	// comparisons or "=" characters inside the value expression as
	// assignment boundaries):
	//
	//	"$x = 'a=b'"      assigns the string "a=b"
	//	"$flag = $a == $b" assigns the comparison's result
	//	"$count += 2"     is a compound assignment
	equals := strings.Index(trimmed, "=")

	if equals > 0 {
		switch trimmed[equals-1] {
		case '=', '!', '<', '>':
			// A comparison operator ("==", "!=", "<=", ">="): the reducer
			// doesn't assign to anything.
		case '+', '-', '*', '/', '%':
			// A compound assignment ("$count += 2") expands into a plain
			// assignment of the compound expression.
			sink = strings.TrimSpace(trimmed[:equals-1])
			if isReactiveSelector(sink) {
				value := strings.TrimSpace(trimmed[equals+1:])

				return sink, "(" + sink + " " + string(trimmed[equals-1]) + " " + value + ")"
			}
		default:
			sink = strings.TrimSpace(trimmed[:equals])
			if isReactiveSelector(sink) {
				return sink, strings.TrimSpace(trimmed[equals+1:])
			}
		}

		return "", ""
	}

	// No "=" at all: the reducer may be a direct function call.
	if isCallReducer(trimmed) {
		return "", trimmed
	}

	return "", ""
}

// isReactiveSelector returns whether the string is a bare reactive variable
// selector. e.g. "$count".
func isReactiveSelector(selector string) bool {
	if len(selector) < 2 || selector[0] != '$' {
		return false
	}

	for i := 1; i < len(selector); i++ {
		if !isIdentifierByte(selector[i]) {
			return false
		}
	}

	return true
}

// isCallReducer returns whether the reducer is a bare function call
// expression, e.g. "save()" or "save($name, $other)".
//
// The call target must be a bare identifier: server script functions are
// imported into the page's shared runtime scope with their own names.
func isCallReducer(reducer string) bool {
	open := strings.Index(reducer, "(")

	if open <= 0 || !strings.HasSuffix(reducer, ")") {
		return false
	}

	target := strings.TrimSpace(reducer[:open])

	if target == "" {
		return false
	}

	for i := 0; i < len(target); i++ {
		if !isIdentifierByte(target[i]) {
			return false
		}
	}

	return true
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
	// The position of the binding in the source file. Errors about the
	// reducer (e.g. a call to a function that wasn't imported from a server
	// script) are reported against it.
	position lexer.Position
	children AbstractSyntaxTree
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
