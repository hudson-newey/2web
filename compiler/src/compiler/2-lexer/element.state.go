package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/2-lexer/states"
)

// When inside the first starting angle bracket (<) and up until (and including)
// the closing angle bracket (>).
func elementLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"!doctype": {token: lexeme.Doctype, next: elementLexer},

		// TODO: Move these out of the element lexer.
		// Having this in here technically means that if you add a script or style
		// attribute to an element, it will switch into the script/style lexer which
		// is not correct.
		// (Or maybe this is a feature? Needs more thought.)
		"script compiled": {token: lexeme.CompiledScriptStartTag, next: inlineCompiledScriptTagLexer},
		"script":          {token: lexeme.ScriptStartTag, next: inlineScriptTagLexer},
		"style":           {token: lexeme.StyleStartTag, next: inlineStyleTagLexer},
		"code":            {token: lexeme.CodeStartTag, next: inlineCodeTagLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, elementLexer)

	return lexerFactory(cases, states.Element)(model)
}
