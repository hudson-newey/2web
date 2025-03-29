package ssg

import (
	"hudson-newey/2web/src/compiler/ssg/modules"
	"strings"
)

type ssgContent struct {
	selector string
	content  string
}

type isStable = bool

func ProcessStaticSite(content string, filePath string) (string, isStable) {
	ssgContent := extractSsgContent(content)
	ssgResult := content

	for _, ssgItem := range ssgContent {
		ssgCommand := strings.Split(ssgItem.content, " ")
		selectorContent := ""

		if len(ssgCommand) >= 2 {
			ssgKeyword := ssgCommand[0]

			switch ssgKeyword {
			case includeToken:
				selectorContent = modules.IncludeSsgContent(ssgCommand[1], filePath)

			case forToken:
				selectorContent = modules.ForSsgContent(ssgCommand[1], ssgCommand[2])
			}
		}

		ssgResult = strings.Replace(ssgResult, ssgItem.selector, selectorContent, 1)
	}

	// by comparing the original content with the result, we can determine if
	// the content is stable
	unstable := strings.Contains(ssgResult, ssgStartToken) && strings.Contains(ssgResult, ssgEndToken)

	return ssgResult, !unstable
}

func extractSsgContent(content string) []ssgContent {
	resultContent := []ssgContent{}

	rawContent := ""
	inSsg := false

	for i := range content {
		if i+len(ssgStartToken) > len(content) {
			break
		}

		if content[i:i+len(ssgStartToken)] == ssgStartToken {
			inSsg = true
			rawContent = ""
			continue
		}

		if content[i:i+len(ssgEndToken)] == ssgEndToken {
			inSsg = false

			// remove any characters that are a part of the start token
			// that were remaining in the raw content
			// and strip any leading/trailing whitespace
			rawContent = rawContent[len(ssgStartToken)-1:]
			refinedContent := strings.TrimSpace(rawContent)

			contentObject := ssgContent{
				selector: ssgStartToken + rawContent + ssgEndToken,
				content:  refinedContent,
			}

			resultContent = append(resultContent, contentObject)
			continue
		}

		if inSsg {
			rawContent += string(content[i])
		}
	}

	return resultContent
}
