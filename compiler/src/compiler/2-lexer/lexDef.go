package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"maps"
	"slices"
	"strings"
)

type LexMatcher string

type V2Lexer func(l *Lexer) (V2LexNode, LexFunc)

type lexDef struct {
	token lexerTokens.LexToken
	next  LexFunc
}

type lexDefMap map[LexMatcher]lexDef

func (lexMap *lexDefMap) matching(lexerModel *Lexer, state states.LexState) (V2LexNode, LexFunc) {
	// We use the same logic for multi-character matchers as we do for single
	// character matchers because the peek logic is the same.
	dereferencedMap := *lexMap

	// Using the golang "range" over a map is non-deterministic, meaning that the
	// order of iteration can change from one execution to the next. This can lead
	// to inconsistent behavior and bugs that are hard to trace because they only
	// appear some of the time.
	// By iterating over the map in a sorted order, we ensure that the output of
	// the program is consistent and predictable, meaning that the compiler is
	// more reliable.
	for _, matcher := range slices.Sorted(maps.Keys(dereferencedMap)) {
		definition := dereferencedMap[matcher]

		stringifiedMatcher := string(matcher)

		matcherOffset := len(stringifiedMatcher)
		peeked := lexerModel.peek(matcherOffset)

		// All of my tokens are case-insensitive and whitespace agnostic.
		// Therefore, I can use EqualFold here to compare the two strings without
		// worrying about case sensitivity.
		// E.g. !DOCTYPE and !doctype tokens are equivalent.
		if strings.EqualFold(peeked, stringifiedMatcher) {
			lexerModel.skip(matcherOffset)

			position := Position{
				Row: lexerModel.Pos.Row,
				Col: lexerModel.Pos.Col - matcherOffset,
			}

			lexNode := V2LexNode{
				Pos:     position,
				Token:   definition.token,
				State:   state,
				Content: stringifiedMatcher,
			}

			return lexNode, definition.next
		}
	}

	return V2LexNode{}, nil
}

// Merges two lexDefMaps together, with the src map taking precedence
// over the dst map in the event of key collisions.
func (dst lexDefMap) with(src lexDefMap) lexDefMap {
	maps.Copy(dst, src)
	return dst
}
