package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

// The lexer for when the <script> tag has been opened and before the first >
// meaning that we are technically still in an element tag.
func inlineScriptTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		">": {token: lexerTokens.ScriptEndTag, next: scriptContentLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, inlineScriptTagLexer)

	return lexerFactory(cases, states.ScriptSource)(model)
}

func scriptContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"</script>": {token: lexerTokens.ScriptEndTag, next: textLexer},
	}

	return lexerFactory(cases, states.ScriptSource)(model)
}
