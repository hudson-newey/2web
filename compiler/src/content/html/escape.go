package html

import "strings"

var escapeReplacer = strings.NewReplacer(
	"&", "&amp;",
	"\"", "&quot;",
	"'", "&#x27;",
	"<", "&lt;",
	">", "&gt;",
)

func EscapeHtml(htmlFragment string) string {
	return escapeReplacer.Replace(htmlFragment)
}
