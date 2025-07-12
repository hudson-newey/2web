package lexer

import (
	"strings"
)

var attributeEndToken LexerSelector = []string{">", " ", "\n", "/"}

// A specialized lexer designed and optimized to extract attribute nodes from
// html content.
func FindPropNodes[T voidNode](
	content string,
	prefix LexerSelector,
) []LexNode[T] {
	resultContent := []LexNode[T]{}

	queryPrefix := prefix[0]

	nodeBufferContent := ""
	inNode := false
	inQuote := false

	codeBlockStart := []string{"<code"}
	codeBlockEnd := []string{"</code>"}
	inCodeBlock := false

	for i := range content {
		if !inNode && content[i:i+len(queryPrefix)] == queryPrefix {
			inNode = true
			nodeBufferContent = ""
			continue
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

func trimLast(s string) string {
	// Edge case: return as-is if too short
	if len(s) < 2 {
		return s
	}
	return s[:len(s)-1]
}
