package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var controlFlowFor = Grammar{
	// @for (item of $myArray) {
	// 	<p>{{ $item }}</p>
	// }
	//
	// Whitespace around the loop source and body is optional, so both of
	// these forms are valid:
	//
	//	@for(item of $myArray){<p>{{ $item }}</p>}
	//	@for (item of $myArray) {
	//		<p>{{ $item }}</p>
	//	}
	Def: newDefinition(
		lexeme.AtSymbol,
		lexeme.ControlFlowForKeyword,

		// Whitespace between the keyword and the loop source
		optional(lexeme.TextContent),

		lexeme.BracketOpen,

		// The loop source. e.g. "item of $myArray"
		// This is captured instead of matched exactly so that whitespace
		// around the "of" keyword is supported.
		lexeme.NewCaptureUntil(),
		lexeme.BracketClosed,

		// Whitespace between the loop source and the body
		optional(lexeme.TextContent),

		lexeme.CurlyOpen,

		// The loop body. e.g. "<p>{{ $item }}</p>"
		// Text outputs inside the body survive the capture because the lexer
		// matches a doubled closing brace as one token.
		lexeme.NewCaptureUntil(),
		lexeme.CurlyClosed,
	),
	Constructor: wrapConstructor(nodes.NewControlFlowForNode),
}
