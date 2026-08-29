package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

// The lexer for when the <script> tag has been opened and before the first >
// meaning that we are technically still in an element tag.
func inlineScriptTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		">": {token: lexeme.GreaterAngle, next: scriptContentLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, inlineScriptTagLexer)

	return lexerFactory(cases, scriptSource)(model)
}

func scriptContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"</script>": {token: lexeme.ScriptEndTag, next: textLexer},
	}

	return lexerFactory(cases, scriptSource)(model)
}
