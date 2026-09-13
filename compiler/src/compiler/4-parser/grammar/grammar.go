package grammar

import (
	"slices"

	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
	"strings"
)

type Grammar struct {
	// A sequence of tokens that define the reactive variable
	Def definition

	// A constructor function to create a node from the matched tokens. The
	// parse context lets constructors report syntax errors with the position
	// that caused them and degrade gracefully instead of panicking.
	Constructor func(lexNodes []*lexer.V2LexNode, context *nodes.ParseContext) *nodes.Node

	// Any child grammar definitions that should be recursively applied within
	// this grammar once matched.
	ChildDefs []Grammar
}

// MinimumTokenCount returns the smallest number of input tokens that this
// grammar can possibly match against.
//
// Optional tokens can be missing from the input, so they are excluded. The
// parser uses this as its look-ahead bound.
func (model *Grammar) MinimumTokenCount() int {
	return model.Def.minimumTokenCount()
}

// Matches lexer nodes against the given grammar.
// Returns the matched subset.
func (model *Grammar) Match(lexNodes []*lexer.V2LexNode, context *nodes.ParseContext) []*lexer.V2LexNode { // If we have not processed enough tokens to have a match yet, we can quickly
	// return false.
	//
	// Optional tokens can be absent from the input, so the minimum number of
	// input tokens needed for a match is the number of required (non-optional)
	// tokens in the definition.
	if len(lexNodes) < model.Def.minimumTokenCount() {
		return []*lexer.V2LexNode{}
	}

	// Because a grammar definition can sometimes come after multiple non-matching
	// tokens, we want to check only the last N tokens where N is the length of
	// the grammar definition.
	matchingSubset := lexNodes
	// matchingSubset := lexNodes[len(lexNodes)-len(model.Def):]

	matchedSubset := []*lexer.V2LexNode{}

	// Use an outer index token so that it can increment independently of the
	// loop iteration (e.g. so we can increment the index in the captureUntil
	// blocks).
	index := 0
	// The definition position is tracked explicitly. The capture block below
	// needs the NEXT definition entry (the capture's break condition), which
	// must not be confused with the current position in the input tokens.
	for defIndex, token := range model.Def {
		// If the definition token is optional, we consume the input token when
		// it matches and continue matching the next definition token when it
		// does not.
		if lexeme.IsOptionalToken(token) {
			if index < len(matchingSubset) && matchingSubset[index].Token == lexeme.UnwrapOptional(token) {
				matchedSubset = append(matchedSubset, matchingSubset[index])
				index += 1
			}

			continue
		}

		// If the definition token is an alternation, we consume the input token
		// when it matches any of the alternation's lexemes and fail the match
		// when it does not.
		if lexeme.IsOrToken(token) {
			if index >= len(matchingSubset) || !slices.Contains(lexeme.UnwrapOr(token), matchingSubset[index].Token) {
				return []*lexer.V2LexNode{}
			}

			matchedSubset = append(matchedSubset, matchingSubset[index])
			index += 1

			continue
		}

		// If we come across a CaptureUntil token, we want to continue looping
		// (and incrementing i) until we find the next token (break condition).
		if lexeme.IsSpecialToken(token, lexeme.CaptureUntil) {
			// If the CaptureUntil is the very last token in the definition, the
			// parser would end up in an infinite loop. That is a compiler bug
			// rather than a user error, but panicking would kill the whole
			// build, so the capture block is disabled and the error is reported
			// with the position of the token the match started at.
			//
			// Use i instaed of i+1 because i indexed from zero while the length
			// is indexed from 1.
			if defIndex+1 == len(model.Def) {
				if context != nil {
					context.ReportError(
						"internal grammar error: unbreakable capture block in grammar definition "+model.Def.String(),
						positionOf(lexNodes),
					)
				}

				return []*lexer.V2LexNode{}
			}

			breakTokens := captureBreakTokens(model.Def[defIndex+1])
			for {
				// If we reach the end of the file while searching for the break
				// condition, the capture block was never terminated. Report it
				// as a syntax error (with the position the capture started at)
				// and degrade to the captured content so far.
				if index >= len(matchingSubset) {
					if context != nil {
						context.ReportError(
							"unterminated capture block: the closing token for the block was never found",
							positionOf(lexNodes),
						)
					}

					return matchedSubset
				}

				if slices.Contains(breakTokens, matchingSubset[index].Token) {
					break
				}

				matchedSubset = append(matchedSubset, matchingSubset[index])
				index += 1
			}

			continue
		}

		if index >= len(matchingSubset) || matchingSubset[index].Token != token {
			return []*lexer.V2LexNode{}
		}

		matchedSubset = append(matchedSubset, matchingSubset[index])
		index += 1
	}

	return matchedSubset
}

// captureBreakTokens returns the input tokens that end a CaptureUntil block.
//
// The break condition can itself be an optional or alternation token, in which
// case the capture block must end on the token(s) that it wraps.
func captureBreakTokens(breakCondition lexeme.Lexeme) []lexeme.Lexeme {
	if lexeme.IsOptionalToken(breakCondition) {
		return []lexeme.Lexeme{lexeme.UnwrapOptional(breakCondition)}
	}

	if lexeme.IsOrToken(breakCondition) {
		return lexeme.UnwrapOr(breakCondition)
	}

	return []lexeme.Lexeme{breakCondition}
}

// positionOf returns the position of the first token in the slice.
func positionOf(lexNodes []*lexer.V2LexNode) lexer.Position {
	if len(lexNodes) == 0 {
		return lexer.StartingPosition
	}

	return lexNodes[0].Pos
}

func (model *definition) String() string {
	var tokens []string
	for _, token := range *model {
		tokens = append(tokens, token.String())
	}
	return "[" + strings.Join(tokens, ", ") + "]"
}
