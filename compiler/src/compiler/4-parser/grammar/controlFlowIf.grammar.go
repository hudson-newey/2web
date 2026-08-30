package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var controlFlowIf = Grammar{
	// @if (expression) {
	//	<div>Hello World!</div>
	// }
	Def: newDefinition(
		lexeme.AtSymbol,
		lexeme.ControlFlowIfKeyword,
		lexeme.BracketOpen,
		lexeme.TextContent,
		lexeme.BracketClosed,
		lexeme.CurlyOpen,
		// Conditional content to show
		lexeme.TextContent,
		lexeme.CurlyClosed,
	),
	// We have to include TextRules as a sub-grammar so that control flow like
	// if conditions can use things like property references and also contain
	// sub-conditions like an if inside of an if.
	// ChildDefs:   []Grammar{TextRules},
	Constructor: wrapConstructor(nodes.NewControlFlowIfNode),
}
