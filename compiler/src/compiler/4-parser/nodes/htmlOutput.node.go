package nodes

import (
	"fmt"
	"regexp"
	"strings"

	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
	"hudson-newey/2web/src/models"
)

// NewHtmlOutputNode compiles an html output:
//
//	[[ $htmlContent ]]
//	[[ getUserTable({ searchQuery: "hello" }) ]]
//
// Unlike a text output ('{{ $content }}'), which renders its expression's
// value as text (it lowers into a *textContent property binding), an html
// output renders the value as html (it lowers into an *innerHTML property
// binding). This makes it possible to render html content that a reactive
// variable holds, or html that a remote function returns over rpc.
//
// The output compiles into one of two flavors:
//
//   - A reactive expression (e.g. '$htmlContent') lowers into an innerHTML
//     property binding, so the page re-renders when the value changes.
//   - An expression that calls a function (e.g. a remote function over rpc)
//     can't be assigned synchronously (the call is asynchronous), so the
//     container is rendered by an async load block: the function is called
//     once when the page loads and its result is assigned to the container.
//     Reactive variables that the expression reads are passed with the value
//     they had at load time.
func NewHtmlOutputNode(lexNodes []*lexer.V2LexNode, context *ParseContext) Node {
	expression, err := scanners.CapturedContent(lexNodes, lexeme.DoubleSquareOpen, lexeme.DoubleSquareClosed)
	if err != nil {
		return context.DegradedNode(
			"html output is missing an expression. Html outputs are written with '[[ expression ]]'",
			lexNodes,
		)
	}

	expression = strings.TrimSpace(expression)

	if expression == "" {
		return context.DegradedNode(
			"html output is missing an expression. Html outputs are written with '[[ expression ]]'",
			lexNodes,
		)
	}

	if containsFunctionCall(expression) {
		// The container's compiler selector is a unique placeholder that the
		// compilation pass replaces with the runtime selector (the same
		// mechanism that reactive property bindings use).
		placeholder := fmt.Sprintf("data-__2_html_%d", context.nextSyntheticId())

		return &htmlOutputNode{
			expression:      expression,
			containerMarkup: fmt.Sprintf(`<span %s></span>`, placeholder),
			placeholder:     placeholder,
			async:           true,
		}
	}

	return &htmlOutputNode{
		expression: expression,
		reactiveProp: &reactivePropertyNode{
			propName: "innerHTML",
			reducer:  expression,
			markupContent: fmt.Sprintf(
				`<span *innerHTML="%s"></span>`,
				expression,
			),
		},
	}
}

// containsFunctionCall returns whether the expression calls a function.
//
// A call makes the expression's value asynchronous (every remote function
// returns a promise), so the value can't be assigned by the synchronous
// property machinery (see NewHtmlOutputNode).
func containsFunctionCall(expression string) bool {
	callPattern := regexp.MustCompile(`[A-Za-z_$][A-Za-z0-9_$.]*\s*\(`)

	return callPattern.MatchString(expression)
}

type htmlOutputNode struct {
	expression string

	// The reactive flavor lowers into an innerHTML property binding (the same
	// structure that a text output lowers into, but with innerHTML instead of
	// textContent).
	reactiveProp *reactivePropertyNode

	// The async flavor renders into a container element that this node owns.
	// containerMarkup is the container's raw markup (with a unique compiler
	// placeholder as its selector) and is replaced during compilation.
	containerMarkup string
	placeholder     string
	async           bool

	children AbstractSyntaxTree
}

func (m *htmlOutputNode) Type() string {
	return "htmlOutputNode"
}

func (m *htmlOutputNode) Children() AbstractSyntaxTree {
	if m.reactiveProp != nil {
		return AbstractSyntaxTree{m.reactiveProp}
	}

	return m.children
}

func (m *htmlOutputNode) MarkupContent() string {
	if m.reactiveProp != nil {
		// The container is owned by the property node (see Children), which is
		// rendered during the markup pass.
		return ""
	}

	return m.containerMarkup
}

func (m *htmlOutputNode) Content(pageModel *page.Page, index *ReactiveIndex) NodeContent {
	if !m.async {
		m.compileReactiveOutput(pageModel, index)
	} else {
		m.compileAsyncOutput(pageModel, index)
	}

	return NodeContent{
		HtmlContent:      pageModel.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

// compileReactiveOutput validates that the expression actually reads a
// reactive variable. An expression without one can never be wired up (the
// property machinery compiles property bindings through the variables they
// read), so the output degrades to its literal source text instead of
// emitting a broken runtime assignment.
func (m *htmlOutputNode) compileReactiveOutput(pageModel *page.Page, index *ReactiveIndex) {
	if len(variableSelectorsReferencedBy(m.expression, index.Variables)) > 0 {
		return
	}

	literalError := models.NewError(
		fmt.Sprintf(
			"html output '[[ %s ]]' does not read any reactive variable. Html outputs render reactive variables (e.g. '[[ $htmlContent ]]') or function calls (e.g. '[[ renderTable({ query: \"hi\" }) ]]')",
			m.expression,
		),
		pageModel.InputPath,
		lexer.Position{},
	)

	documentErrors.AddErrors(&literalError)

	// Degrade the output to its literal source so that the author's content
	// stays visible in the page.
	pageModel.SetContent(
		strings.ReplaceAll(
			pageModel.Html.Content,
			m.reactiveProp.markupContent,
			"[[ "+m.expression+" ]]",
		),
	)
}

// compileAsyncOutput wires the container and emits the async load block that
// calls the expression's function and renders its result.
func (m *htmlOutputNode) compileAsyncOutput(pageModel *page.Page, index *ReactiveIndex) {
	elementSelector := pageModel.Ids.CreateElementName()
	pageModel.SetContent(
		strings.ReplaceAll(pageModel.Html.Content, m.placeholder, elementSelector),
	)

	// Reactive variables that the expression reads are passed with the value
	// they have when the load block runs (runtime variables by reference, and
	// static variables inlined with their value).
	expression := resolveReactiveExpression(m.expression, index)

	loadBlock := fmt.Sprintf(
		`(async () => {
	try {
		document.querySelector("[%s]")["innerHTML"] = "" + (await %s);
	} catch (__2_error) {
		console.error("[2web] failed to render html output '%s':", __2_error);
	}
})();`,
		elementSelector, expression, strings.ReplaceAll(m.expression, `"`, `\"`),
	)

	pageModel.AppendReactiveRuntime(page.ReactiveRuntimeChunk{
		Variable:     "html:" + elementSelector,
		Dependencies: variableSelectorsReferencedBy(m.expression, index.Variables),
		Content:      loadBlock,
	})
}

func (m *htmlOutputNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *htmlOutputNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}
