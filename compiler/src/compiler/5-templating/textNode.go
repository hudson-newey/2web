package templating

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/5-templating/reactiveCompiler"
	"strings"
)

var textStartToken lexer.LexerSelector = []string{"{{"}
var textEndToken lexer.LexerSelector = []string{"}}"}

func expandTextNodes(content string) string {
	textNodes := lexer.FindNodes[lexer.TextNode](content, textStartToken, textEndToken)

	for _, node := range textNodes {
		// text nodes can have reducers if they modify the content of the
		// inner reactive variable
		// E.g. {{ $count * 2 }}
		nodeReducer := strings.Join(node.Tokens[0:len(node.Tokens)], " ")

		doubleQuotes := reactiveCompiler.UseDoubleQuotes(nodeReducer)

		// While adding DOM nodes may decrease the performance of the
		// application without further processing, the compiler can correctly
		// tree shake away this element if it's inner text is never changed
		// by a reactive event.
		// Meaning that the node will be removed if targeted DOM updated are not
		// needed.
		// Additionally, if the innerText of the new element will be updated,
		// it is likely to be more efficient because we don't have to re-render
		// all of the innerText up to the nearest ancestor. We can instead
		// just update specific parts/words/section(s) of the element.
		//
		// e.g. In the example <p>Hello <span *innerText="$message"></span></p>
		// we don't have to re-render the <p> tag of the "Hello" text every
		// time, we can just re-render the "$message" part of the paragraph.
		domNodeElement := ""
		if doubleQuotes {
			domNodeElement = fmt.Sprintf("<span *innerText=\"%s\"></span>", nodeReducer)
		} else {
			domNodeElement = fmt.Sprintf("<span *innerText='%s'></span>", nodeReducer)
		}

		content = strings.ReplaceAll(content, node.Selector, domNodeElement)
	}

	return content
}
