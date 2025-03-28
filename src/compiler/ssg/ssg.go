package ssg

import (
	"strings"
)

type ssgContent struct {
	selector string
	content  string
}

func ProcessStaticSite(content string, filePath string) string {
	ssgContent := extractSsgContent(content)

	for _, ssgItem := range ssgContent {
		ssgCommand := strings.Split(ssgItem.content, " ")
		selectorContent := ""

		if len(ssgCommand) >= 2 {
			ssgKeyword := ssgCommand[0]

			switch ssgKeyword {
			case includeToken:
				selectorContent = includeSsgContent(ssgCommand[1], filePath)
			}
		}

		println(ssgItem.selector, selectorContent)
		content = strings.Replace(content, ssgItem.selector, selectorContent, 1)
	}

	return content
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
