package html

import "strings"

func EscapeHtml(htmlFragment string) string {
	sr := strings.NewReplacer(
		"&", "&amp;",
		"\"", "&quot;",
		"'", "&#x27;",
		"<", "&lt;",
		">", "&gt;",
	)

	return sr.Replace(htmlFragment)
}
