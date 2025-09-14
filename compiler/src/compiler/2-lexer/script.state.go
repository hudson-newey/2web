package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

func scriptLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"</script>": {token: lexerTokens.ScriptEndTag, next: textLexer},
	}

	cases = withStrings(cases, scriptLexer)

	return lexerFactory(cases, states.ScriptSource)(model)
}
