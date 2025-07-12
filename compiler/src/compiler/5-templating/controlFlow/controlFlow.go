package controlFlow

import (
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/5-templating/controlFlow/cfModules/cfFor"
	"hudson-newey/2web/src/compiler/5-templating/controlFlow/cfModules/cfIf"
	"strings"
)

// ssg start and end tokens are exported because we use the same syntax inside
// control flow such as if conditions e.g. {% if <condition> <result> %}
// I have chosen to use the same if start and end tokens for ssg and control
// flow because one of the goals of this project is to make the difference
// between ssg, ssr, isr, and csr an implementation detail that the user
// doesn't have to deal with.
// The compiler should be able to automatically pick the most efficient
// rendering method depending on the circumstances.
var cfStartSelector lexer.LexerSelector = []string{"{%"}
var cfEndSelector lexer.LexerSelector = []string{"%}"}

func ProcessControlFlow(filePath string, content string) string {
	ssgContent := lexer.FindNodes[lexer.SsgNode](content, cfStartSelector, cfEndSelector)
	result := content

	for _, node := range ssgContent {
		selectorContent := ""

		if len(node.Tokens) >= 2 {
			ssgKeyword := node.Tokens[0]

			// we "continue" in the default case so that unrecognized keywords are
			// preserved so that ssg might be able to process them later on
			//
			// TODO: this is currently flawed because control flow inside of ssg
			// "include" statements will not work
			// we should therefore have a unified method to run both control flow and
			// ssg at compile time until the template is stable.
			switch ssgKeyword {
			case forToken[0]:
				selectorContent = cfFor.ForLoopContent(node.Tokens[1], node.Tokens[2:len(node.Tokens)])
			case ifToken[0]:
				selectorContent = cfIf.IfConditionContent(node.Tokens[1], node.Tokens[2:len(node.Tokens)])
			default:
				continue
			}
		}

		result = strings.Replace(result, node.Selector, selectorContent, 1)
	}

	return result
}
