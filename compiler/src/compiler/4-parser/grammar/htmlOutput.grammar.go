package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var htmlOutput = Grammar{
	// [[ $htmlContent ]]
	// [[ getUserTable({ searchQuery: "hello" }) ]]
	//
	// The expression is captured (instead of matched exactly) so that
	// expressions made of multiple tokens (e.g. function calls with object
	// literal arguments) are supported.
	Def: newDefinition(
		lexeme.DoubleSquareOpen,
		lexeme.NewCaptureUntil(),
		lexeme.DoubleSquareClosed,
	),
	Constructor: wrapConstructor(nodes.NewHtmlOutputNode),
}
