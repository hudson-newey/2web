package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
)

// The lexer for when the <script> tag has been opened and before the first >
// meaning that we are technically still in an element tag.
func inlineCompiledScriptTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		">": {token: lexerTokens.GreaterAngle, next: compiledScriptContentLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, compiledScriptContentLexer)

	return lexerFactory(cases, states.CompiledScriptSource)(model)
}

// We want to allow the user to exit a "<script compiled>" block even if they
// have incorrecetly formatted scripts.
// e.g. They forgot a semi colon.
// We do this so that we can parse as much of the content and continue searching
// for more errors in the document unrelated to the script.
func withScriptExitCase(src lexDefMap) lexDefMap {
	exitCases := lexDefMap{
		"</script>": {token: lexerTokens.ScriptEndTag, next: textLexer},
	}
	return src.with(exitCases)
}

func compiledScriptContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"$": {token: lexerTokens.DollarSign, next: reactiveVarAssignmentLexer},
	}
	cases = withScriptExitCase(cases)

	return lexerFactory(cases, states.CompiledScriptSource)(model)
}

func reactiveVarAssignmentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"=": {token: lexerTokens.Equals, next: reactiveVarAssignmentLexer},
		";": {token: lexerTokens.Semicolon, next: compiledScriptContentLexer},
	}
	cases = withScriptExitCase(cases)
	return lexerFactory(cases, states.CompiledScriptSource)(model)
}
