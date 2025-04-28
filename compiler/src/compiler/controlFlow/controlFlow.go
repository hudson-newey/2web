package controlFlow

import (
	"hudson-newey/2web/src/compiler/controlFlow/cfModules/cfFor"
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

			switch ssgKeyword {
			case forToken[0]:
				selectorContent = cfFor.ForLoopContent(node.Tokens[1], node.Tokens[2])
			default:
				continue
			}
		}

		result = strings.Replace(result, node.Selector, selectorContent, 1)
	}

	return result
}
