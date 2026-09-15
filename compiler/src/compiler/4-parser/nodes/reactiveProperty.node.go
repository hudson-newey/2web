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

func NewReactivePropertyNode(lexNodes []*lexer.V2LexNode, context *ParseContext) Node {
	propName, err := scanners.NthToken(lexNodes, lexeme.TextContent, 1)
	if err != nil {
		return context.DegradedNode(
			"reactive property binding is missing a property name. Properties are bound with '*property=\"reducer\"'",
			lexNodes,
		)
	}

	reducer, err := scanners.NthToken(lexNodes, lexeme.TextContent, 2)
	if err != nil {
		return context.DegradedNode(
			"reactive property binding is missing a reducer. Properties are bound with '*property=\"reducer\"'",
			lexNodes,
		)
	}

	markupContent := fmt.Sprintf(
		"*%s=\"%s",
		propName.Content,
		reducer.Content,
	)

	return &reactivePropertyNode{
		propName:      strings.TrimSpace(propName.Content),
		reducer:       strings.TrimSpace(reducer.Content),
		markupContent: markupContent,
	}
}

type reactivePropertyNode struct {
	propName      string
	reducer       string
	markupContent string
	children      AbstractSyntaxTree

	// optional:
	// If the property node has already been used during the compilation process
	// it will already have a selector allocated to it that we can reuse.
	allocatedSelector string

	// negated inverts the value that is assigned to the property at runtime.
	//
	// This is needed for props like `hidden`, where `@if (condition)` should
	// show its content when the condition is truthy, but the `hidden` property
	// hides an element when it is assigned a truthy value.
	negated bool
}

// propValue wraps the given runtime value expression in the property's value
// transformation (e.g. negation for the `hidden` property).
func (m *reactivePropertyNode) propValue(value string) string {
	if m.negated {
		return fmt.Sprintf("!(%s)", value)
	}

	return value
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

func (m *reactivePropertyNode) Content(page *page.Page, _ *ReactiveIndex) NodeContent {
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

func (m *reactivePropertyNode) selector(pageModel *page.Page) string {
	if m.allocatedSelector != "" {
		return m.allocatedSelector
	}

	// Replace the compile time selector with a runtime selector
	// TODO: Use definitions from lexer here instead
	compilerSelector := fmt.Sprintf("*%s=\"%s\"", m.propName, m.reducer)
	domSelector := pageModel.Ids.CreateElementName()

	m.allocatedSelector = domSelector

	pageModel.SetContent(
		strings.Replace(pageModel.Html.Content, compilerSelector, domSelector, 1),
	)

	return domSelector
}

// Some properties such as text assignment can be evaluated at compile time
// meaning that we don't need to ship any runtime code.
func (m *reactivePropertyNode) canCompilerInline() bool {
	supportedPropSinks := []string{"textContent", "innerText"}
	return slices.Contains(supportedPropSinks, m.propName)
}

// propAssignment wraps the runtime value expression for the property's
// assignment.
//
// String sinks (textContent/innerText) stringify the assigned value. DOM
// stringification of a falsy value (e.g. the number 0) behaves differently
// between engines, so the emitted code makes the string conversion explicit.
func (m *reactivePropertyNode) propAssignment(value string) string {
	if m.canCompilerInline() {
		return fmt.Sprintf(`"" + (%s)`, value)
	}

	return value
}

// Finds all reactive variable dependencies for a property binding
func (m *reactivePropertyNode) reactiveVariableDeps(ast AbstractSyntaxTree) []*reactiveVariableNode {
	return lists.Filter(ast.reactiveVariables(), func(x *reactiveVariableNode) bool {
		return strings.Contains(m.reducer, x.selector())
	})
}
