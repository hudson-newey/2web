package nodes

import (
	"fmt"
	"hash/fnv"
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

// NewControlFlowForNode compiles a reactive for loop:
//
//	@for (xyz of $myArray) {
//		<p>index value: {{ $xyz }}</p>
//	}
//
// The loop is lowered into two synthesized reactive nodes so that the whole
// existing reactivity machinery (dependency tracking, update cascades, and
// initial rendering) applies to it:
//
//   - A computed variable whose initial value maps the loop's source array to
//     the rendered body markup (one bound per item, with the loop variable
//     bound to the item):
//
//     $__2_for_<id> = $myArray.map(function ($xyz) { return `...body...`; }).join("")
//
//   - An innerHTML property binding on a container element that renders the
//     computed variable's value.
//
// When the source array is reassigned (by an event or a compiled function),
// the computed variable is re-evaluated by the array's update cascade and the
// container's innerHTML is updated.
func NewControlFlowForNode(lexNodes []*lexer.V2LexNode, context *ParseContext) Node {
	loopSource, err := scanners.CapturedContent(lexNodes, lexeme.BracketOpen, lexeme.BracketClosed)
	if err != nil {
		return context.DegradedNode(
			"@for block is missing a loop source. For blocks are declared with '@for (item of $myArray) { content }'",
			lexNodes,
		)
	}

	body, err := scanners.CapturedContent(lexNodes, lexeme.CurlyOpen, lexeme.CurlyClosed)
	if err != nil {
		return context.DegradedNode(
			"@for block is missing a body. For blocks are declared with '@for (item of $myArray) { content }'",
			lexNodes,
		)
	}

	loopVar, sourceSelector, err := parseForLoopSource(loopSource)
	if err != nil {
		return context.DegradedNode(
			err.Error()+" For blocks are declared with '@for (item of $myArray) { content }'",
			lexNodes,
		)
	}

	// The synthesized variable's name is derived from the loop's source so
	// that the compiled output is deterministic, and is made unique per loop
	// instance so that two identical loops don't declare the same runtime
	// variable twice.
	id := fmt.Sprintf("%08x", fnv32(loopSource+"\x00"+body))[:8]
	instance := context.nextControlFlowId(id)

	variableName := "__2_for_" + id
	if instance > 1 {
		variableName = fmt.Sprintf("%s_%d", variableName, instance)
	}

	initialValue := fmt.Sprintf(
		`%s.map(function ($%s) { return `+"`%s`"+`; }).join("")`,
		sourceSelector, loopVar, compileForLoopBody(body),
	)

	// The loop's rendered markup is owned by the property node (which is
	// rendered during the markup pass), so the for node itself must not
	// render any markup. Rendering both would duplicate the container.
	reactiveProp := &reactivePropertyNode{
		propName: "innerHTML",
		reducer:  "$" + variableName,
		markupContent: fmt.Sprintf(
			`<span *innerHTML="$%s"></span>`,
			variableName,
		),
	}

	return &controlFlowForNode{
		loopVar:          loopVar,
		sourceSelector:   sourceSelector,
		body:             body,
		computedVariable: newSynthesizedVariable(variableName, initialValue),
		reactiveProp:     reactiveProp,
	}
}

// parseForLoopSource parses the loop variable and source array out of the
// loop's captured source. e.g. "xyz of $myArray".
func parseForLoopSource(loopSource string) (loopVar string, sourceSelector string, err error) {
	trimmed := strings.TrimSpace(loopSource)

	sourcePattern := regexp.MustCompile(
		`^([A-Za-z_$][A-Za-z0-9_$]*)\s+of\s+([A-Za-z_$][A-Za-z0-9_$]*)$`,
	)

	match := sourcePattern.FindStringSubmatch(trimmed)
	if match == nil {
		return "", "", fmt.Errorf(
			"@for loop source must be in the form 'item of $myArray', got '%s'.",
			trimmed,
		)
	}

	loopVar = strings.TrimPrefix(match[1], "$")
	sourceSelector = match[2]

	if !strings.HasPrefix(sourceSelector, "$") {
		return "", "", fmt.Errorf(
			"@for loop source must be a reactive variable (it is missing the $ prefix), got '%s'.",
			sourceSelector,
		)
	}

	return loopVar, sourceSelector, nil
}

// compileForLoopBody converts the loop's body markup into the template that
// the generated map callback returns for every item:
//
//   - Text outputs ('{{ $item }}') become template literal interpolations
//     ('${ $item }'), so that every reactive expression the body supports in
//     plain markup also works inside a loop body.
//   - The body is escaped for a template literal, so that quotes, backticks,
//     and interpolations in the user's markup can't break the generated code.
func compileForLoopBody(body string) string {
	escaped := escapeTemplateLiteral(body)

	textOutputPattern := regexp.MustCompile(`\{\{\s*(.*?)\s*\}\}`)

	return textOutputPattern.ReplaceAllString(escaped, "${ $1 }")
}

// escapeTemplateLiteral escapes a markup fragment so that it can be embedded
// in a JavaScript template literal.
func escapeTemplateLiteral(markup string) string {
	escaped := strings.ReplaceAll(markup, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, "`", "\\`")
	escaped = strings.ReplaceAll(escaped, "${", `$\{`)

	return escaped
}

// fnv32 hashes the given text into a 32 bit fnv1a hash.
func fnv32(text string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(text))

	return hash.Sum32()
}

// newSynthesizedVariable creates a reactive variable node that wasn't written
// by the user (see controlFlowForNode).
func newSynthesizedVariable(variableName string, initialValue string) *reactiveVariableNode {
	return &reactiveVariableNode{
		variableName: variableName,
		initialValue: initialValue,
	}
}

type controlFlowForNode struct {
	// The name the loop variable is bound to inside the body. e.g. "xyz"
	loopVar string

	// The selector of the reactive variable that the loop iterates.
	// e.g. "$myArray"
	sourceSelector string

	// The raw body markup.
	body string

	computedVariable *reactiveVariableNode
	reactiveProp     *reactivePropertyNode
	children         AbstractSyntaxTree
}

func (m *controlFlowForNode) Type() string {
	return "controlFlowForNode"
}

func (m *controlFlowForNode) Children() AbstractSyntaxTree {
	return AbstractSyntaxTree{m.reactiveProp, m.computedVariable}
}

func (m *controlFlowForNode) MarkupContent() string {
	// The loop's container element is owned by the innerHTML property node
	// (see Children), which is rendered during the markup pass.
	return ""
}

func (m *controlFlowForNode) Content(pageModel *page.Page, index *ReactiveIndex) NodeContent {
	if index.VariableBySelector(m.sourceSelector) == nil {
		sourceError := models.NewError(
			fmt.Sprintf(
				"@for loop source '%s' is not a declared reactive variable. Loop sources are declared in a <script compiled> block (e.g. '$myArray = [\"a\", \"b\"];')",
				m.sourceSelector,
			),
			pageModel.InputPath,
			lexer.Position{},
		)

		documentErrors.AddErrors(&sourceError)
	}

	if index.VariableBySelector("$"+m.loopVar) != nil {
		shadowError := models.NewError(
			fmt.Sprintf(
				"the @for loop variable '$%s' shadows a reactive variable that is declared in this page. Rename the loop variable (e.g. '@for (%sItem of %s)')",
				m.loopVar, m.loopVar, m.sourceSelector,
			),
			pageModel.InputPath,
			lexer.Position{},
		)

		documentErrors.AddErrors(&shadowError)
	}

	return NodeContent{
		HtmlContent:      pageModel.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

func (m *controlFlowForNode) AddChild(child Node) {
	m.reactiveProp.AddChild(child)
}

func (m *controlFlowForNode) RemoveChild(child Node) {
	m.reactiveProp.RemoveChild(child)
}
