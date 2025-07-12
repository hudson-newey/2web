package preprocessor

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"strings"
)

var elementRefToken lexer.LexerSelector = []string{"#"}

// Element ref's are literally just aliases for an "id" attribute that is
// tracked by the reactive compiler.
// Warning: This means that element ref's cannot be used if an element has an
// "id" attribute.
func expandElementRefs(content string) string {
	refNodes := lexer.FindPropNodes[lexer.RefNode](content, elementRefToken)

	for _, node := range refNodes {
		// Remove the first character from the node selector because it will be the
		// hash (#) symbol. The id that we want to produce does not include this
		// character.
		idAttribute := fmt.Sprintf(`id="%s"`, node.Selector[1:])
		content = strings.ReplaceAll(content, node.Selector, idAttribute)
	}

	return content
}
