package ssg

import (
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/ssg/modules"
	"strings"
)

type isStable = bool

// TODO: remove the token[0] hacks below
func ProcessStaticSite(filePath string, content string) (string, isStable) {
	ssgContent := lexer.FindNodes[lexer.SsgNode](content, ssgStartToken, ssgEndToken)
	ssgResult := content

	for _, node := range ssgContent {
		selectorContent := ""

		if len(node.Tokens) >= 2 {
			ssgKeyword := node.Tokens[0]

			switch ssgKeyword {
			case includeToken[0]:
				selectorContent = modules.IncludeSsgContent(node.Tokens[1], filePath)
			case forToken[0]:
				selectorContent = modules.ForSsgContent(node.Tokens[1], node.Tokens[2])
			}
		}

		ssgResult = strings.Replace(ssgResult, node.Selector, selectorContent, 1)
	}

	// by comparing the original content with the result, we can determine if
	// the content is stable
	unstable := strings.Contains(ssgResult, ssgStartToken[0]) && strings.Contains(ssgResult, ssgEndToken[0])

	return ssgResult, !unstable
}
