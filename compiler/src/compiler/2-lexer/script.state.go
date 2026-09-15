package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

var inlineScriptTagLexerState stateLexers

// The lexer for when the <script> tag has been opened and before the first >
// meaning that we are technically still in an element tag.
func inlineScriptTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := inlineScriptTagLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			">": {token: lexeme.GreaterAngle, next: scriptContentLexer},
		}

		cases = withAttributes(cases, inlineScriptTagLexer, "inline-script-tag")
		cases = withStrings(cases, inlineScriptTagLexer, "inline-script-tag")
		return cases
	})

	return lexerFactory(matchers, tagAttributes)(model)
}

var scriptContentLexerState stateLexers

func scriptContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := scriptContentLexerState.get(func() lexDefMap {
		return lexDefMap{
			"</script>": {token: lexeme.ScriptEndTag, next: textLexer},
		}
	})

	return lexerFactory(matchers, scriptSource)(model)
}
