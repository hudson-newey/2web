package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

// The lexer for when the <script> tag has been opened and before the first >
// meaning that we are technically still in an element tag.
func inlineCompiledScriptTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		">": {token: lexeme.GreaterAngle, next: compiledScriptContentLexer},
	}

	cases = withAttributes(cases)
	cases = withStrings(cases, compiledScriptContentLexer)

	return lexerFactory(cases, compiledScriptSource)(model)
}

// We want to allow the user to exit a "<script compiled>" block even if they
// have incorrecetly formatted scripts.
// e.g. They forgot a semi colon.
// We do this so that we can parse as much of the content and continue searching
// for more errors in the document unrelated to the script.
func withScriptExitCase(src lexDefMap) lexDefMap {
	exitCases := lexDefMap{
		"</script>": {token: lexeme.ScriptEndTag, next: textLexer},
	}
	return src.with(exitCases)
}

func compiledScriptContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		// process comments first so that commenting out a code block can
		// disable the code (the commented out code never gets lexed output).
		"//":     {token: lexeme.MarkupCommentStart, next: esmLineCommentLexer},
		"/*":     {token: lexeme.ScriptBlockCommentStart, next: esmBlockCommentLexer},
		"$":      {token: lexeme.DollarSign, next: reactiveVarAssignmentLexer},
		"import": {token: lexeme.KeywordImport, next: esmImportLexer},
	}
	cases = withScriptExitCase(cases)
	return lexerFactory(cases, compiledScriptSource)(model)
}

func reactiveVarAssignmentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"=": {token: lexeme.Equals, next: reactiveVarAssignmentLexer},
		";": {token: lexeme.Semicolon, next: compiledScriptContentLexer},
	}
	cases = withScriptExitCase(cases)
	return lexerFactory(cases, compiledScriptSource)(model)
}

func esmImportLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"from": {token: lexeme.KeywordFrom, next: esmImportLexer},
		";":    {token: lexeme.Semicolon, next: compiledScriptContentLexer},
	}
	cases = withScriptExitCase(cases)
	return lexerFactory(cases, compiledScriptSource)(model)
}

func esmLineCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"\n": {token: lexeme.BlockCommentEnd, next: compiledScriptContentLexer},
	}
	cases = withScriptExitCase(cases)
	return lexerFactory(cases, compiledScriptSource)(model)
}

func esmBlockCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	cases := lexDefMap{
		"*/": {token: lexeme.BlockCommentEnd, next: compiledScriptContentLexer},
	}
	cases = withScriptExitCase(cases)
	return lexerFactory(cases, compiledScriptSource)(model)
}
