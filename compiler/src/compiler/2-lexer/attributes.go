package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
)

// Attribute lexing is shared between multiple states.
// e.g. The inline style and script tags can have attribute but require their
// own lexer because any text state they transition into is actually a new lexer
// state where you are lexing script/css text rather than normal html text.
//
// returningState is the lexer that attribute punctuation (space, equals, etc.)
// returns to. Callers that represent a specific tag (e.g. <code ...>) must pass
// themselves so that the tag's own exit conditions (e.g. the closing >) stay
// reachable. Returning to a fixed state (e.g. elementLexer) would lose the tag
// context and the tag would be terminated by the wrong lexer.
//
// cacheKey identifies the calling state for the string lexer cache.
func withAttributes(src lexDefMap, returningState LexFunc, cacheKey string) lexDefMap {
	attributeStates := lexDefMap{
		// I treat tabs like spaces so that they are treated the same in attributes
		" ":  {token: lexeme.Space, next: returningState},
		"\t": {token: lexeme.Space, next: returningState},
		"/":  {token: lexeme.Slash, next: returningState},

		// Quote characters are deliberately NOT registered here. They are
		// registered by withStrings, which transitions into a string lexer state
		// so that tag matchers can't fire inside quoted attribute values.
		//
		// Registering them here would shadow the string lexer (lexDefMap.with
		// doesn't override existing keys) and attribute values would be lexed in
		// element state. In element state, a quoted value like
		// data-testid="html-code" would fire the "code" tag matcher in the
		// middle of the attribute value, splitting the tag in two.
		"#": {token: lexeme.Hash, next: returningState},
		"=": {token: lexeme.Equals, next: returningState},
		"!": {token: lexeme.Exclamation, next: textLexer},
		">": {token: lexeme.GreaterAngle, next: textLexer},

		"*": {token: lexeme.Star, next: reactivePropertyLexer},
		"@": {token: lexeme.AtSymbol, next: reactiveEventLexer},
	}

	return src.with(attributeStates)
}

var reactivePropertyLexerState stateLexers

func reactivePropertyLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := reactivePropertyLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			"=": {token: lexeme.Equals, next: reactivePropertyLexer},
		}
		return withStrings(cases, elementLexer, "reactive-property")
	})

	return lexerFactory(matchers, element)(model)
}

var reactiveEventLexerState stateLexers

func reactiveEventLexer(model *Lexer) (V2LexNode, LexFunc) {
	matchers := reactiveEventLexerState.get(func() lexDefMap {
		cases := lexDefMap{
			"=": {token: lexeme.Equals, next: reactiveEventLexer},
		}
		return withStrings(cases, elementLexer, "reactive-event")
	})

	return lexerFactory(matchers, element)(model)
}
