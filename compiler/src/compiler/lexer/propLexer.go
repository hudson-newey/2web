package lexer

import (
	"strings"
)

var attributeEndToken LexerToken = []string{">", " ", "\n", "/"}

// A specialized lexer designed and optimized to extract attribute nodes from
// html content.
func FindPropNodes[T voidNode](
	content string,
	prefix LexerPropPrefix,
) []LexNode[T] {
	htmlContent := isolateHtmlContent(content)

	resultContent := []LexNode[T]{}

	queryPrefix := prefix[0]

	nodeBufferContent := ""
	inNode := false
	inQuote := false

	for i := range htmlContent {
		if !inNode && htmlContent[i:i+len(queryPrefix)] == queryPrefix {
			inNode = true
			nodeBufferContent = ""
			continue
		}

		if inNode {
			if htmlContent[i] == '"' {
				inQuote = !inQuote
			}

			if !inQuote {
				for _, endToken := range attributeEndToken {
					if htmlContent[i:i+len(endToken)] == endToken {
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
			nodeBufferContent += string(htmlContent[i])
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

func isolateHtmlContent(content string) string {
	nonHtmlStartTags := []string{"<script", "<style", "<math"}
	nonHtmlEndTags := []string{"</script>", "</style>", "</math>"}

	inHtml := true
	resultContent := ""

	for i := range content {
		if !inHtml {
			for _, endTag := range nonHtmlEndTags {
				if i-len(endTag) < 0 {
					continue
				}

				if content[i-len(endTag):i] == endTag {
					inHtml = true
					break
				}
			}
		}

		for _, startTag := range nonHtmlStartTags {
			if i+len(startTag) > len(content) {
				continue
			}

			if content[i:i+len(startTag)] == startTag {
				inHtml = false
			}
		}

		if inHtml {
			resultContent += string(content[i])
		}
	}

	return resultContent
}
