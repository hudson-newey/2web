package grammar

import (
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/ast"
	"strings"
)

type grammar struct {
	// A sequence of tokens that define the reactive variable
	Def definition

	// A constructor function to create a node from the tokens
	Constructor func(lexNodes []*lexer.V2LexNode) *ast.Node

	// Any child grammar definitions that should be recursively applied within
	// this grammar once matched.
	ChildDefs []grammar
}

func (model *grammar) Matches(lexNodes []*lexer.V2LexNode) bool {
	// If we have not processed enough tokens to have a match yet, we can quickly
	// return false.
	if len(lexNodes) < len(model.Def) {
		return false
	}

	// Because a grammar definition can sometimes come after multiple non-matching
	// tokens, we want to check only the last N tokens where N is the length of
	// the grammar definition.
	matchingSubset := lexNodes
	// matchingSubset := lexNodes[len(lexNodes)-len(model.Def):]

	// Use an outer index token so that it can increment independently of the
	// loop iteration (e.g. so we can increment the index in the captureUntil
	// blocks).
	index := 0
	for _, token := range model.Def {
		// If we come across a CaptureUntil token, we want to continue looping
		// (and incrementing i) until we find the next token (break condition).
		if lexerTokens.IsSpecialToken(token, lexerTokens.CaptureUntil) {
			// If the CaptureUntil is the very last token in the definition,
			// the parser would end up in an infinite loop.
			// It's ok to use a panic here since if we panic here, the program
			// would never even run once propperly.
			//
			// Use i instaed of i+1 because i indexed from zero while the length
			// is indexed from 1.
			if index+1 == len(model.Def) {
				panic("Detected unbreakable CaptureUntil block (infinite parser loop)")
			}

			breakCondition := model.Def[index+1]
			for matchingSubset[index].Token != breakCondition {
				// We can skip the first itteration since we know that it'll be
				// the initial capture token.
				index += 1

				// If we reach the end of the file while searching for the break
				// condition, we probably want to log an error.
				if index == len(matchingSubset) {
					cli.PrintWarning("Parser finished inside of capture block")
					return true
				}
			}
			continue
		}

		if matchingSubset[index].Token != token {
			return false
		}

		index += 1
	}

	return true
}

func (model *definition) String() string {
	var tokens []string
	for _, token := range *model {
		tokens = append(tokens, token.String())
	}
	return "[" + strings.Join(tokens, ", ") + "]"
}
