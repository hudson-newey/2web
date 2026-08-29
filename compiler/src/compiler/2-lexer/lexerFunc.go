package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"io"
)

type LexFunc func(*Lexer) (V2LexNode, LexFunc)

func lexerFactory(lexMap lexDefMap, state lexState) LexFunc {
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

			panic(err)
		}

		startPos := lexerModel.Pos
		lexerModel.backup(1)
		text := lexerModel.lexLiteral(lexMap)

		// There are different types of text depending on what context we are in
		// Sometimes it can be external source code.
		tokenMap := map[lexState]lexeme.Lexeme{
			compiledScriptSource: lexeme.CompiledScriptSource,
			scriptSource:         lexeme.ScriptSource,
			styleSource:          lexeme.StyleSource,
			codeSource:           lexeme.CodeSource,
			textContent:          lexeme.TextContent,
		}

		token, exists := tokenMap[state]
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
