package ssg

import (
	"hudson-newey/2web/src/compiler/lexer"
	"hudson-newey/2web/src/compiler/ssg/ssgModules"
	"hudson-newey/2web/src/content/html"
	"strings"
)

type isStable = bool

// TODO: remove the token[0] hacks below
func ProcessStaticSite(filePath string, content string) (string, isStable) {
	ssgContent := lexer.FindNodes[lexer.SsgNode](content, SsgStartToken, SsgEndToken)
	ssgResult := content

	for _, node := range ssgContent {
		selectorContent := ""

		if len(node.Tokens) >= 2 {
			ssgKeyword := node.Tokens[0]

			switch ssgKeyword {
			case includeToken[0]:
				selectorContent = ssgModules.IncludeSsgContent(node.Tokens[1], filePath)
			}
		}

		ssgResult = strings.Replace(ssgResult, node.Selector, selectorContent, 1)
	}

	if html.IsHtmlFile(filePath) {
		ssgResult = ssgModules.ExpandElementRefs(ssgResult)
	}

	// by comparing the original content with the result, we can determine if
	// the content is stable
	unstable := strings.Contains(ssgResult, SsgStartToken[0]) && strings.Contains(ssgResult, SsgEndToken[0])

	return ssgResult, !unstable
}
