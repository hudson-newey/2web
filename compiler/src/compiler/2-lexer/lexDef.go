package lexer

import (
	"hudson-newey/2web/src/compiler/2-lexer/states"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"slices"
	"strings"
)

type LexMatcher = string

type V2Lexer func(l *Lexer) (V2LexNode, LexFunc)

type lexDef struct {
	token lexerTokens.LexToken
	next  LexFunc
}

type lexDefMap map[LexMatcher]lexDef

func (lexMap *lexDefMap) matching(lexerModel *Lexer, state states.LexState) (V2LexNode, LexFunc) {
	dereferencedMap := *lexMap

	// Collect keys into a slice explicitly to avoid relying on maps.Keys which
	// may return an iterator type (not a slice) in newer Go versions.
	var keys []string
	for k := range dereferencedMap {
		keys = append(keys, k)
	}

	// Using the golang "range" over a map is non-deterministic, meaning that the
	// order of iteration can change from one execution to the next. This can lead
	// to inconsistent behavior and bugs that are hard to trace because they only
	// appear some of the time.
	// By iterating over the map in a sorted order, we ensure that the output of
	// the program is consistent and predictable, meaning that the compiler is
	// more reliable.
	//
	// Maximal munch: prefer the longest possible matcher so multi-character
	// tokens (e.g. "<!--") are not shadowed by shorter prefix tokens (e.g. "<").
	//
	// TODO: Review if we want to use maximal munch because it does have some
	// trade-offs in certain scenarios.
	slices.SortFunc(keys, func(a, b string) int {
		if len(a) != len(b) {
			return len(b) - len(a)
		}
		if a == b {
			return 0
		}
		if a < b {
			return -1
		}
		return 1
	})

	for _, matcher := range keys {
		definition := dereferencedMap[matcher]

		matcherOffset := len(matcher)
		peeked := lexerModel.peek(matcherOffset)

		// All of my tokens are case-insensitive and whitespace agnostic.
		// Therefore, I can use EqualFold here to compare the two strings without
		// worrying about case sensitivity.
		// E.g. !DOCTYPE and !doctype tokens are equivalent.
		if strings.EqualFold(peeked, matcher) {
			lexerModel.skip(matcherOffset)

			position := Position{
				Row: lexerModel.Pos.Row,
				Col: lexerModel.Pos.Col - matcherOffset,
			}

			lexNode := V2LexNode{
				Pos:     position,
				Token:   definition.token,
				State:   state,
				Content: matcher,
			}

			return lexNode, definition.next
		}
	}

	return NewV2LexNode(), nil
}

// wouldMatchAt returns true if any matcher in the map matches at the current
// reader position without consuming any input. This is used by lexLiteral to
// detect an exit condition before consuming the character that begins it.
func (lexMap *lexDefMap) wouldMatchAt(lexerModel *Lexer) bool {
	dereferencedMap := *lexMap
	// We can short-circuit on the first match; ordering doesn't affect the
	// boolean outcome. Still, prefer checking longer matchers first to avoid
	// unnecessary shorter peeks.
	var keys []string
	for k := range dereferencedMap {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b string) int {
		if len(a) != len(b) {
			return len(b) - len(a)
		}
		if a == b {
			return 0
		}
		if a < b {
			return -1
		}
		return 1
	})

	for _, matcher := range keys {
		if strings.EqualFold(lexerModel.peek(len(matcher)), matcher) {
			return true
		}
	}
	return false
}

// Merges two lexDefMaps together, with the dst map taking precedence
// over the src map in the event of key collisions.
func (dst lexDefMap) with(src lexDefMap) lexDefMap {
	for k, v := range src {
		if _, exists := dst[k]; !exists {
			dst[k] = v
		}
	}

	return dst
}
