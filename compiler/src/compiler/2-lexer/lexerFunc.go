package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"io"
)

type LexFunc func(*Lexer) (V2LexNode, LexFunc)

// contentTokens maps a lexer state to the token that literal text captured in
// that state is emitted as. The map is constant, so it is built once instead
// of per lexer invocation.
var contentTokens = map[lexState]lexeme.Lexeme{
	compiledScriptSource: lexeme.CompiledScriptSource,
	scriptSource:         lexeme.ScriptSource,
	styleSource:          lexeme.StyleSource,
	codeSource:           lexeme.CodeSource,
	textContent:          lexeme.TextContent,
}

func lexerFactory(lexMap *compiledMatchers, state lexState) LexFunc {
	return func(lexerModel *Lexer) (V2LexNode, LexFunc) {
		matchingLexNode, nextState := lexMap.matching(lexerModel, state)
		if nextState != nil {
			return matchingLexNode, nextState
		}

		// Some systems use \r\n for new lines, but because we ignore \r, it will
		// work fine.
		readerChar, _, err := lexerModel.nextChar()
		if readerChar == '\n' {
			lexerModel.lineFeed()

			// We do not want to include new lines in the final lexed output, so we
			// keep lexing until we find a different token.
			for {
				return lexerFactory(lexMap, state)(lexerModel)
			}
		}

		if err != nil {
			if err == io.EOF {
				position := Position{
					Row: lexerModel.Pos.Row,
					Col: lexerModel.Pos.Col,
				}

				// We still return a textLexer state here so that if a file is partially
				// corrupted and contains an EOF before the real end, we can recover
				// gracefully.
				lexNode := V2LexNode{
					Pos:     position,
					Token:   lexeme.EOF,
					Content: "",
					State:   sourceText,
				}

				return lexNode, lexerModel.State
			}

			// A read failure (other than end of file) must not kill the build.
			lexerModel.RecordError("failed to read source: "+err.Error(), *lexerModel.Pos)
			lexNode := V2LexNode{
				Pos:     *lexerModel.Pos,
				Token:   lexeme.EOF,
				State:   sourceText,
				Content: "",
			}

			return lexNode, lexerModel.State
		}

		// Copy the position by value. Pos is a pointer, so capturing it by
		// reference would alias the live position: every mutation made while
		// scanning the literal (and by backup below) would move the recorded
		// start along with it, making every content token report the position
		// where the literal ENDED instead of where it started.
		startPos := *lexerModel.Pos
		lexerModel.backup(1)
		text := lexerModel.lexLiteral(lexMap)

		// There are different types of text depending on what context we are in
		// Sometimes it can be external source code.
		token, exists := contentTokens[state]
		if !exists {
			token = lexeme.TextContent
		}

		position := Position{
			Row: startPos.Row,
			Col: startPos.Col,
		}

		lexNode := V2LexNode{
			Pos:     position,
			Token:   token,
			State:   state,
			Content: string(text),
		}

		return lexNode, lexerModel.State
	}
}
