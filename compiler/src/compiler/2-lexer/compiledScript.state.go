package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

var inlineCompiledScriptTagLexerState stateLexers

// The lexer for when the <script> tag has been opened and before the first >
// meaning that we are technically still in an element tag.
func inlineCompiledScriptTagLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := inlineCompiledScriptTagLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			">": {token: lexeme.GreaterAngle, next: compiledScriptContentLexer},
		}

		cases = withAttributes(cases, inlineCompiledScriptTagLexer, "inline-compiled-script-tag")
		cases = withStrings(cases, inlineCompiledScriptTagLexer, "inline-compiled-script-tag")
		return cases
	})

	return lexerFactory(matchers, tagAttributes)(model)
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

var compiledScriptContentLexerState stateLexers

func compiledScriptContentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := compiledScriptContentLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			// process comments first so that commenting out a code block can
			// disable the code (the commented out code never gets lexed output).
			"//":     {token: lexeme.MarkupCommentStart, next: esmLineCommentLexer},
			"/*":     {token: lexeme.ScriptBlockCommentStart, next: esmBlockCommentLexer},
			"$":      {token: lexeme.DollarSign, next: reactiveVarAssignmentLexer},
			"import": {token: lexeme.KeywordImport, next: esmImportLexer},
		}
		cases = withScriptExitCase(cases)
		return cases
	})

	return lexerFactory(matchers, compiledScriptSource)(model)
}

var reactiveVarAssignmentLexerState stateLexers

func reactiveVarAssignmentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := reactiveVarAssignmentLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			"=": {token: lexeme.Equals, next: reactiveVarAssignmentLexer},
			";": {token: lexeme.Semicolon, next: compiledScriptContentLexer},
		}
		cases = withScriptExitCase(cases)
		return cases
	})

	return lexerFactory(matchers, compiledScriptSource)(model)
}

var esmImportLexerState stateLexers

func esmImportLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := esmImportLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			"from": {token: lexeme.KeywordFrom, next: esmImportLexer},
			";":    {token: lexeme.Semicolon, next: compiledScriptContentLexer},
		}
		cases = withScriptExitCase(cases)
		return cases
	})

	return lexerFactory(matchers, compiledScriptSource)(model)
}

var esmLineCommentLexerState stateLexers

func esmLineCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := esmLineCommentLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			"\n": {token: lexeme.BlockCommentEnd, next: compiledScriptContentLexer},
		}
		cases = withScriptExitCase(cases)
		return cases
	})

	return lexerFactory(matchers, compiledScriptSource)(model)
}

var esmBlockCommentLexerState stateLexers

func esmBlockCommentLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := esmBlockCommentLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			"*/": {token: lexeme.BlockCommentEnd, next: compiledScriptContentLexer},
		}
		cases = withScriptExitCase(cases)
		return cases
	})

	return lexerFactory(matchers, compiledScriptSource)(model)
}
