package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"io"
	"strings"
)

var textLexerState stateLexers

func textLexer(model *Lexer) (V2LexNode, LexFunc) {
	// Escape sequences are handled before the generic matchers so that the
	// escape character itself is never emitted into the page. The backslash
	// must stay in the matchers below so that literal captures (lexLiteral)
	// stop at escape characters instead of swallowing them.
	if strings.EqualFold(model.peek(1), "\\") {
		model.skip(1)
		return escapeSequenceLexer(model)
	}

	matchers := textLexerState.get(func() lexDefMap {
		return lexDefMap{
			// The comment state has to always come first that it takes precedence over
			// other matches and can omit them as source text.
			"<!--": {token: lexeme.MarkupCommentStart, next: markupCommentLexer},

			"<":  {token: lexeme.LessAngle, next: elementLexer},
			">":  {token: lexeme.GreaterAngle, next: textLexer},
			"\\": {token: lexeme.Escape, next: escapeSequenceLexer},

			"{": {token: lexeme.CurlyOpen, next: textLexer},
			"}": {token: lexeme.CurlyClosed, next: textLexer},

			// The html output delimiters. e.g. '[[ $htmlContent ]]'
			"[[": {token: lexeme.DoubleSquareOpen, next: textLexer},
			"]]": {token: lexeme.DoubleSquareClosed, next: textLexer},

			"@": {token: lexeme.AtSymbol, next: controlFlowLexer},
		}
	})

	return lexerFactory(matchers, textContent)(model)
}

// escapeSequenceLexer consumes the character that follows an escape
// character (e.g. the ">" of "\>") and emits it as literal text content.
//
// The escaped character is emitted as a text content token (instead of the
// token it would normally lex as) so that the parser can't interpret it. e.g.
// an escaped "\>" must not be lexed as the end of an element, and an escaped
// "\{" must not be able to form a reactive text output.
func escapeSequenceLexer(model *Lexer) (V2LexNode, LexFunc) {
	position := Position{
		Row: model.Pos.Row,
		Col: model.Pos.Col,
	}

	char, _, err := model.nextChar()
	model.Pos.Col++

	if err != nil {
		if err == io.EOF {
			// A trailing escape character at the end of the file has nothing to
			// escape; recover gracefully by returning to the text state.
			return V2LexNode{
				Pos:     position,
				Token:   lexeme.EOF,
				State:   sourceText,
				Content: "",
			}, textLexer
		}

		// A read failure (other than end of file) must not kill the build.
		model.RecordError("failed to read source: "+err.Error(), *model.Pos)
		return V2LexNode{
			Pos:     position,
			Token:   lexeme.EOF,
			State:   sourceText,
			Content: "",
		}, textLexer
	}

	if char == '\n' {
		model.lineFeed()
	}

	return V2LexNode{
		Pos:     position,
		Token:   lexeme.TextContent,
		State:   textContent,
		Content: string(char),
	}, textLexer
}
