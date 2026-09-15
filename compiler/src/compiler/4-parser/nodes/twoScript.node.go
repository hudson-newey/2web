package nodes

import (
	"strings"

	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewTwoScriptNode(lexNodes []*lexer.V2LexNode, context *ParseContext) Node {
	// Find the lexNode that is a StyleSource token
	var content string
	for _, lexNode := range lexNodes {
		if lexNode.Token == lexeme.CompiledScriptSource {
			content = lexNode.Content
			break
		}
	}

	return &twoScriptNode{
		lexerNodes: lexNodes,
		content:    content,
	}
}

type twoScriptNode struct {
	lexerNodes []*lexer.V2LexNode
	content    string
	children   AbstractSyntaxTree

	// compiledFunctions caches the top level function declarations that were
	// extracted from the block source (nil until the first extraction).
	compiledFunctions []compiledFunction
}

func (m *twoScriptNode) Type() string {
	return "twoScriptNode"
}

func (m *twoScriptNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *twoScriptNode) MarkupContent() string {
	return ""
}

// functions extracts the top level function declarations from the block
// source. The result is memoized so that the index build and the compilation
// pass share one extraction.
func (m *twoScriptNode) functions() []compiledFunction {
	if m.compiledFunctions == nil {
		m.compiledFunctions = extractCompiledFunctions(m.lexerNodes)
	}

	return m.compiledFunctions
}

// prepareCompiledFunctions reconciles the block's parsed children with its
// function declarations and must run before the page's reactive index is
// built:
//
//   - The lexer has no brace context, so a '$name = value;' statement inside a
//     function body is parsed exactly like a top level reactive variable
//     declaration. Those statements belong to the function, not to the page's
//     reactive state, so they are removed from the AST.
//   - A variable that is only ever assigned inside a function body has no
//     declaration at all. One is synthesized (with an empty initial value) so
//     that it gets a runtime representation and can be rendered.
//
// It returns the block's function declarations.
func (m *twoScriptNode) prepareCompiledFunctions() []compiledFunction {
	functions := m.functions()

	kept := make(AbstractSyntaxTree, 0, len(m.children))
	declared := map[string]bool{}

	for _, child := range m.children {
		variable, isVariable := child.(*reactiveVariableNode)
		if isVariable && m.isInsideFunction(variable.position) {
			// The declaration is an assignment inside a function body.
			continue
		}

		if isVariable {
			declared[variable.selector()] = true
		}

		kept = append(kept, child)
	}

	for _, function := range functions {
		for _, reference := range function.references {
			if reference.kind == refRead || declared[reference.selector] {
				continue
			}

			// Synthesize the variable's declaration with an empty initial
			// value: its real (first) value only exists once the function
			// has run.
			declared[reference.selector] = true
			kept = append(kept, &reactiveVariableNode{
				variableName: strings.TrimPrefix(reference.selector, "$"),
				initialValue: `""`,
				position: lexer.Position{
					Row: reference.position.row,
					Col: reference.position.col,
				},
			})
		}
	}

	m.children = kept

	return functions
}

// isInsideFunction returns whether the given source position falls within one
// of the block's function declarations.
func (m *twoScriptNode) isInsideFunction(position lexer.Position) bool {
	for _, function := range m.functions() {
		if positionWithinFunction(
			sourcePosition{row: position.Row, col: position.Col},
			function,
		) {
			return true
		}
	}

	return false
}

func (m *twoScriptNode) Content(pageModel *page.Page, index *ReactiveIndex) NodeContent {
	// The block's functions are emitted into the page's shared reactive
	// runtime (the same scope as the reactive variables and the server rpc
	// passthroughs) with their reactive references rewritten.
	for _, function := range m.functions() {
		if function.name == "" {
			continue
		}

		pageModel.AppendReactiveRuntime(page.ReactiveRuntimeChunk{
			Variable: "fn:" + function.name,
			Content:  rewriteCompiledFunction(function, index),
		})
	}

	return NodeContent{
		HtmlContent:      pageModel.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.FromContent(m.content),
	}
}

func (m *twoScriptNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *twoScriptNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

// positionWithinFunction returns whether the position falls within the
// function's declaration span (inclusive of the keyword and the closing
// brace).
func positionWithinFunction(position sourcePosition, function compiledFunction) bool {
	return position.row >= function.start.row && position.row <= function.end.row &&
		!(position.row == function.start.row && position.col < function.start.col) &&
		!(position.row == function.end.row && position.col > function.end.col)
}
