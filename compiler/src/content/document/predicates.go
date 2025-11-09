package document

import "strings"

func HasBodyTag(content string) bool {
	return strings.Contains(strings.ToLower(content), "<body")
}
