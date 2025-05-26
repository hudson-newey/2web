package html

import "strings"

func EscapeHtml(htmlFragment string) string {
	replacementTable := map[string]string{
		"{": "&#x7B;",
		"}": "&#x7D;",
		"<": "&lt;",
		">": "&gt;",
	}

	result := htmlFragment
	for key, value := range replacementTable {
		result = strings.ReplaceAll(result, key, value)
	}

	return result
}
