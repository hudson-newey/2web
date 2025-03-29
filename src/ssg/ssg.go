package ssg

import (
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/ssg/modules"
	"strings"
)

type isStable = bool

func ProcessStaticSite(filePath string, content string) (string, isStable) {
	ssgContent := lexer.FindNodes[lexer.SsgNode](content, ssgStartToken, ssgEndToken)
	ssgResult := content

	for _, node := range ssgContent {
		selectorContent := ""

		if len(node.Tokens) >= 2 {
			ssgKeyword := node.Tokens[0]

			switch ssgKeyword {
			case includeToken:
				selectorContent = modules.IncludeSsgContent(node.Tokens[1], filePath)
			case forToken:
				selectorContent = modules.ForSsgContent(node.Tokens[1], node.Tokens[2])
			}
		}

		ssgResult = strings.Replace(ssgResult, node.Selector, selectorContent, 1)
	}

	// by comparing the original content with the result, we can determine if
	// the content is stable
	unstable := strings.Contains(ssgResult, ssgStartToken) && strings.Contains(ssgResult, ssgEndToken)

	return ssgResult, !unstable
}
