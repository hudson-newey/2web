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

			contentObject := lexNode{
				Selector: startToken + rawContent + endToken,
				Content:  refinedContent,
				Tokens:   refinedSplitTokens,
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
