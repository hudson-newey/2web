package lexer

import "strings"

type lexNode struct {
	Selector string
	Content  string
	Tokens   []string
}

func FindNodes(content string, startToken string, endToken string) []lexNode {
	resultContent := []lexNode{}

	rawContent := ""
	inNode := false

	for i := range content {
		if i+len(startToken) > len(content) {
			break
		}

		if content[i:i+len(startToken)] == startToken {
			inNode = true
			rawContent = ""
			continue
		}

		if content[i:i+len(endToken)] == endToken {
			inNode = false

			// remove any characters that are a part of the start token
			// that were remaining in the raw content
			// and strip any leading/trailing whitespace
			rawContent = rawContent[len(startToken)-1:]
			refinedContent := strings.TrimSpace(rawContent)
			splitTokens := strings.Split(refinedContent, " ")

			contentObject := lexNode{
				Selector: startToken + rawContent + endToken,
				Content:  refinedContent,
				Tokens:   splitTokens,
			}

			resultContent = append(resultContent, contentObject)
			continue
		}

		if inNode {
			rawContent += string(content[i])
		}
	}

	return resultContent
}
