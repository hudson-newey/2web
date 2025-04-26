package lexer

import (
	"strings"
)

// TODO: this should not be exported from this module
var attributeEndToken LexerToken = []string{">", " ", "\n", "/"}

func FindPropNodes[T voidNode](
	content string,
	prefix LexerPropPrefix,
) []LexNode[T] {
	resultContent := []LexNode[T]{}

	queryPrefix := prefix[0]

	nodeBufferContent := ""
	inNode := false
	inQuote := false

	for i := range content {
		if !inNode && content[i:i+len(queryPrefix)] == queryPrefix {
			inNode = true
			nodeBufferContent = ""
			continue
		}

		if inNode {
			if content[i] == '"' {
				inQuote = !inQuote
			}

			if !inQuote {
				for _, endToken := range attributeEndToken {
					if content[i:i+len(endToken)] == endToken {
						refinedContent := strings.TrimSpace(nodeBufferContent)

						// remember that attributes look like
						// @click="$count = $count + 1"
						// in this case, the tokens would be
						// []string{"click", "$count", "=", "$count", "+", "1"}
						// notice that the equals sign is not represented, and everything
						// between the first and last quotation marks are split

						// notice that we only replace the first equals sign
						// this means that the value can contain equal signs
						tokenContents := strings.Replace(refinedContent, "=", " ", 1)

						// replace the first quotation mark
						tokenContents = strings.Replace(tokenContents, "\"", "", 1)

						// remove the  last character (the closing quotation mark)
						tokenContents = trimLast(tokenContents)

						nodeBufferTokens := strings.Split(tokenContents, " ")

						contentObject := LexNode[T]{
							Selector: queryPrefix + nodeBufferContent,
							Content:  refinedContent,
							Tokens:   nodeBufferTokens,
						}

						resultContent = append(resultContent, contentObject)
						inNode = false
						break
					}
				}

			}
		}

		if inNode {
			nodeBufferContent += string(content[i])
		}
	}

	return resultContent
}

func FindNodes[T voidNode](
	content string,
	candidateStartTokens LexerToken,
	candidateEndTokens LexerToken,
) []LexNode[T] {
	resultContent := []LexNode[T]{}

	startToken := candidateStartTokens[0]
	endToken := candidateEndTokens[0]

	nodeBufferContent := ""
	inNode := false

	for i := range content {
		if i+len(startToken) > len(content) {
			break
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

				token = strings.TrimPrefix(token, "\"")
				token = strings.TrimPrefix(token, "'")

				token = strings.TrimSuffix(token, "\"")
				token = strings.TrimSuffix(token, "'")

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

func trimLast(s string) string {
	// Edge case: return as-is if too short
	if len(s) < 2 {
		return s
	}
	return s[:len(s)-1]
}
