package lexer

import (
	"bytes"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"slices"
	"strings"
	"sync"
)

type LexMatcher = string

type V2Lexer func(l *Lexer) (V2LexNode, LexFunc)

type lexDef struct {
	token lexeme.Lexeme
	next  LexFunc
}

type lexDefMap map[LexMatcher]lexDef

// matcherEntry is a single matcher with its definition, pre-lowercased so that
// the case insensitive comparison can run over byte slices without allocating.
type matcherEntry struct {
	matcher LexMatcher
	folded  []byte
	def     lexDef
}

// compiledMatchers is a lexDefMap that has been flattened into a sorted
// matcher slice.
//
// The lexers previously rebuilt (and re-sorted) their matcher maps on every
// invocation, which dominated compile time profiles. The matcher sets are
// constant per lexer state, so they are built once and reused.
//
// Entries are sorted by matcher length (descending) so that maximal munch
// matching is preserved: multi character tokens (e.g. "<!--") are checked
// before shorter prefix tokens (e.g. "<").
type compiledMatchers struct {
	entries []matcherEntry
}

// compile flattens a lexDefMap into a compiledMatchers value.
func compileMatchers(defs lexDefMap) *compiledMatchers {
	entries := make([]matcherEntry, 0, len(defs))

	maxLength := 0
	for matcher, def := range defs {
		entries = append(entries, matcherEntry{
			matcher: matcher,
			folded:  []byte(strings.ToLower(matcher)),
			def:     def,
		})

		if len(matcher) > maxLength {
			maxLength = len(matcher)
		}
	}

	// Maximal munch: prefer the longest possible matcher so multi-character
	// tokens (e.g. "<!--") are not shadowed by shorter prefix tokens (e.g. "<").
	slices.SortFunc(entries, func(a, b matcherEntry) int {
		if len(a.matcher) != len(b.matcher) {
			return len(b.matcher) - len(a.matcher)
		}
		return strings.Compare(a.matcher, b.matcher)
	})

	return &compiledMatchers{entries: entries}
}

// stateLexers lazily compiles the matcher set for a single lexer state.
//
// Every state's matcher set is constant, so each state owns one of these and
// builds its set on first use.
type stateLexers struct {
	once     sync.Once
	matchers *compiledMatchers
}

// get compiles the matcher set on first use and returns it.
//
// The build function is only invoked once; it typically composes a state's
// base matchers with withAttributes and withStrings.
func (s *stateLexers) get(build func() lexDefMap) *compiledMatchers {
	s.once.Do(func() {
		s.matchers = compileMatchers(build())
	})

	return s.matchers
}

// matching finds the first matcher that matches at the current reader
// position and emits its token.
//
// All matchers are case-insensitive (e.g. !DOCTYPE and !doctype tokens are
// equivalent).
func (c *compiledMatchers) matching(lexerModel *Lexer, state lexState) (V2LexNode, LexFunc) {
	for _, entry := range c.entries {
		peeked := lexerModel.peekBytes(len(entry.matcher))

		if bytes.EqualFold(peeked, entry.folded) {
			matcherOffset := len(entry.matcher)
			lexerModel.skip(matcherOffset)

			position := Position{
				Row: lexerModel.Pos.Row,
				Col: lexerModel.Pos.Col - matcherOffset,
			}

			lexNode := V2LexNode{
				Pos:     position,
				Token:   entry.def.token,
				State:   state,
				Content: entry.matcher,
			}

			return lexNode, entry.def.next
		}
	}

	return NewV2LexNode(), nil
}

// wouldMatchAt returns true if any matcher matches at the current reader
// position without consuming any input. This is used by lexLiteral to detect
// an exit condition before consuming the character that begins it.
func (c *compiledMatchers) wouldMatchAt(lexerModel *Lexer) bool {
	for _, entry := range c.entries {
		peeked := lexerModel.peekBytes(len(entry.matcher))

		if bytes.EqualFold(peeked, entry.folded) {
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
