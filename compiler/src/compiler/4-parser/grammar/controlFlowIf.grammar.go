package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var controlFlowIf = Grammar{
	// @if (expression) {
	//	<div>Hello World!</div>
	// }
	//
	// Whitespace around the condition and body is optional, so both of these
	// forms are valid:
	//
	//	@if($isOpen){Hello World!}
	//	@if ($isOpen) {
	//		Hello World!
	//	}
	Def: newDefinition(
		lexeme.AtSymbol,
		lexeme.ControlFlowIfKeyword,

		// Whitespace between the keyword and the condition
		optional(lexeme.TextContent),

		lexeme.BracketOpen,

		// The condition expression. e.g. "$isOpen"
		// This is captured instead of matched exactly so that expressions made
		// of multiple tokens (e.g. "$count > 5") are supported.
		lexeme.NewCaptureUntil(),
		lexeme.BracketClosed,

		// Whitespace between the condition and the body
		optional(lexeme.TextContent),

		lexeme.CurlyOpen,

		// Conditional content to show
		lexeme.NewCaptureUntil(),
		lexeme.CurlyClosed,
	),
	// We have to include TextRules as a sub-grammar so that control flow like
	// if conditions can use things like property references and also contain
	// sub-conditions like an if inside of an if.
	// ChildDefs:   []Grammar{TextRules},
	Constructor: wrapConstructor(nodes.NewControlFlowIfNode),
}
