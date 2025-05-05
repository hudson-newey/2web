package controlFlow

import (
	"hudson-newey/2web/src/compiler/controlFlow/cfModules/cfFor"
	"hudson-newey/2web/src/compiler/controlFlow/cfModules/cfIf"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/ssg"
	"strings"
)

func ProcessControlFlow(filePath string, content string) string {
	ssgContent := lexer.FindNodes[lexer.SsgNode](content, ssg.SsgStartToken, ssg.SsgEndToken)
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
