package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

// When inside the first starting angle bracket (<) and up until (and including)
// the closing angle bracket (>).
func elementLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"!doctype": {token: lexerTokens.Doctype, next: elementLexer},

		// TODO: Move these out of the element lexer.
		// Having this in here technically means that if you add a script or style
		// attribute to an element, it will switch into the script/style lexer which
		// is not correct.
		// (Or maybe this is a feature? Needs more thought.)
		"script compiled": {token: lexerTokens.CompiledScriptStartTag, next: inlineCompiledScriptTagLexer},
		"script":          {token: lexerTokens.ScriptStartTag, next: inlineScriptTagLexer},
		"style":           {token: lexerTokens.StyleStartTag, next: inlineStyleTagLexer},
		"code":            {token: lexerTokens.CodeStartTag, next: inlineCodeTagLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, elementLexer)

	return lexerFactory(cases, states.Element)(model)
}
