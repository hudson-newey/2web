package lexer

import (
	"strings"
)

func FindNodes[T voidNode](
	content string,
	candidateStartTokens LexerSelector,
	candidateEndTokens LexerSelector,
) []LexNode[T] {
	resultContent := []LexNode[T]{}

	startToken := candidateStartTokens[0]
	endToken := candidateEndTokens[0]

	nodeBufferContent := ""
	inNode := false

	codeBlockStart := []string{"<code"}
	codeBlockEnd := []string{"</code>"}
	inCodeBlock := false

	for i := range content {
		if inNode {
			if i+len(endToken) > len(content) {
				break
			}
		} else {
			if i+len(startToken) > len(content) {
				break
			}
		}

		for _, startTag := range codeBlockStart {
			if i+len(startTag) > len(content) {
				continue
			}

			if content[i:i+len(startTag)] == startTag {
				inCodeBlock = true
				break
			}
		}

		for _, endTag := range codeBlockEnd {
			if i-len(endTag) < 0 {
				continue
			}

			if content[i-len(endTag):i] == endTag {
				inCodeBlock = false
				break
			}
		}

		if inCodeBlock {
			continue
		}

		if !inNode && content[i:i+len(startToken)] == startToken {
			inNode = true
			nodeBufferContent = ""
			continue
		}

		if inNode && content[i:i+len(endToken)] == endToken {
			inNode = false

			// remove any characters that are a part of the start token
			// that were remaining in the raw content
			// and strip any leading/trailing whitespace
			nodeBufferContent = nodeBufferContent[len(startToken)-1:]
			refinedContent := strings.TrimSpace(nodeBufferContent)
			splitTokens := strings.Split(refinedContent, " ")

			refinedSplitTokens := []string{}
			tokenBuffer := ""
			isInQuote := false

			// if there is a token in the splitTokens that starts with a
			// quotation mark (single or double), and does not end with
			// a quotation mark, then we want to continue to the next
			// token, joining the tokens together until we find a token that
			// ends with a quotation mark
			for _, token := range splitTokens {
				hasQuotePrefix := strings.HasPrefix(token, "\"") || strings.HasPrefix(token, "'")
				hasQuoteSuffix := strings.HasSuffix(token, "\"") || strings.HasSuffix(token, "'")

				// token = strings.TrimPrefix(token, "\"")
				// token = strings.TrimPrefix(token, "'")

				// token = strings.TrimSuffix(token, "\"")
				// token = strings.TrimSuffix(token, "'")

				if hasQuotePrefix && !hasQuoteSuffix {
					tokenBuffer += token + " "
					isInQuote = true
					continue
				}

				if isInQuote {
					tokenBuffer += token + " "
					if hasQuoteSuffix {
						isInQuote = false
						refinedSplitTokens = append(refinedSplitTokens, strings.TrimSpace(tokenBuffer))
						tokenBuffer = ""
					}
					continue
				}

				refinedSplitTokens = append(refinedSplitTokens, token)
			}

			contentObject := LexNode[T]{
				Selector: startToken + nodeBufferContent + endToken,
				Content:  refinedContent,
				Tokens:   refinedSplitTokens,
			}

			resultContent = append(resultContent, contentObject)
			continue
		}

		if inNode {
			nodeBufferContent += string(content[i])
		}
	}

	return resultContent
}
