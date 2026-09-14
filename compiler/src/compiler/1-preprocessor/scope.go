package preprocessor

import (
	"fmt"
	"regexp"
	"strings"
)

// Component styles are scoped with the css "@scope" at rule.
//
// A component instance's top level elements are tagged with a per instance
// scoping attribute (e.g. data-__2_scope="3"), and every style block of the
// component is wrapped in "@scope ([data-__2_scope="3"]) { ... }", so the
// component's styles only apply to its own markup (and its descendants) and
// can't leak into the rest of the page (or into other instances of the same
// component).
//
// The scope ids are deterministic per page (they increment in document
// order), so repeated builds compile identical output.

// The attribute that marks the scoping root elements of a component instance.
const scopeAttributeFormat = `data-__2_scope="%d"`

// Matches a style block of a component (including its tags).
var styleBlockPattern = regexp.MustCompile(`(?is)(<style[^>]*>)(.*?)(</style>)`)

// hasStyleBlock returns whether the component content declares any styles.
// Components without styles don't need scoping (and don't get the scoping
// attributes, keeping their markup untouched).
func hasStyleBlock(content string) bool {
	return strings.Contains(strings.ToLower(content), "<style")
}

// scopeComponent tags the top level elements of a component instance with the
// instance's scoping attribute and wraps the component's style blocks in
// "@scope" rules for that attribute.
func scopeComponent(content string, scopeId int) string {
	attribute := fmt.Sprintf(scopeAttributeFormat, scopeId)

	content = injectScopeAttribute(content, attribute)

	return styleBlockPattern.ReplaceAllStringFunc(content, func(match string) string {
		groups := styleBlockPattern.FindStringSubmatch(match)

		return groups[1] + "\n@scope ([" + attribute + "]) {\n" + groups[2] + "\n}\n" + groups[3]
	})
}

// injectScopeAttribute adds the given attribute to the start tag of every top
// level element of the content.
//
// The content is scanned instead of parsed, because the component's script
// and style blocks contain characters (comparison operators, quoted strings)
// that would confuse a naive tag scanner. Script and style blocks are copied
// verbatim (script blocks don't render markup, and style blocks are scoped
// separately), so nothing inside them is scanned for tags.
func injectScopeAttribute(content string, attribute string) string {
	var out strings.Builder

	i := 0
	depth := 0

	for i < len(content) {
		rest := content[i:]

		switch {
		case strings.HasPrefix(rest, "<!--"):
			out.WriteString(rest[:4])
			i += 4
			i += copyUntil(&out, content[i:], "-->", true)

		case isBlockStart(rest, "script"):
			tagEnd := findTagEnd(rest)
			out.WriteString(rest[:tagEnd])
			i += tagEnd
			i += copyUntil(&out, content[i:], "</script>", true)

		case isBlockStart(rest, "style"):
			tagEnd := findTagEnd(rest)
			out.WriteString(rest[:tagEnd])
			i += tagEnd
			i += copyUntil(&out, content[i:], "</style>", true)

		case strings.HasPrefix(rest, "</"):
			tagEnd := findTagEnd(rest)
			out.WriteString(rest[:tagEnd])
			i += tagEnd

			if depth > 0 {
				depth--
			}

		case strings.HasPrefix(rest, "<"):
			tagEnd := findTagEnd(rest)
			tag := rest[:tagEnd]

			if depth == 0 {
				out.WriteString(injectAttributeIntoTag(tag, attribute))
			} else {
				out.WriteString(tag)
			}

			i += tagEnd

			if opensElement(tag) {
				depth++
			}

		default:
			out.WriteByte(content[i])
			i++
		}
	}

	return out.String()
}

// isBlockStart returns whether the content starts with the given block
// element (case insensitive), and that it opens a block (it isn't an end tag
// or a self closing tag).
func isBlockStart(content string, name string) bool {
	if len(content) < len(name)+1 {
		return false
	}

	if content[0] != '<' {
		return false
	}

	return strings.EqualFold(content[1:1+len(name)], name)
}

// findTagEnd returns the length of the tag that starts at the beginning of
// the content (respecting quoted attribute values, which may contain ">").
func findTagEnd(content string) int {
	inQuote := byte(0)

	for i := 0; i < len(content); i++ {
		char := content[i]

		if inQuote != 0 {
			if char == inQuote {
				inQuote = 0
			}

			continue
		}

		switch char {
		case '"', '\'':
			inQuote = char
		case '>':
			return i + 1
		}
	}

	return len(content)
}

// tagName returns the (lowercase) tag name of a start or end tag.
func tagName(tag string) string {
	trimmed := strings.TrimLeft(tag, "<")
	trimmed = strings.TrimLeft(trimmed, "/")

	name := ""
	for i := 0; i < len(trimmed); i++ {
		char := trimmed[i]

		if char == ' ' || char == '\t' || char == '\n' || char == '\r' || char == '>' || char == '/' {
			break
		}

		name += string(char)
	}

	return strings.ToLower(name)
}

var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

// opensElement returns whether a start tag opens an element that end tags
// close (void elements and self closing tags don't).
func opensElement(tag string) bool {
	if strings.HasSuffix(tag, "/>") {
		return false
	}

	return !voidElements[tagName(tag)]
}

// injectAttributeIntoTag adds an attribute to a start tag.
func injectAttributeIntoTag(tag string, attribute string) string {
	if !strings.HasSuffix(tag, ">") {
		return tag
	}

	closing := tag[len(tag)-2:]

	if closing == "/>" {
		return tag[:len(tag)-2] + " " + attribute + " />"
	}

	// Keep a single space between the last attribute (or the tag name) and
	// the injected attribute.
	if len(tag) >= 2 && tag[len(tag)-2] == ' ' {
		return tag[:len(tag)-1] + attribute + ">"
	}

	return tag[:len(tag)-1] + " " + attribute + ">"
}

// copyUntil copies the content up to (and optionally including) the first
// occurrence of the given end marker, and returns how many bytes were copied.
func copyUntil(out *strings.Builder, content string, end string, includeEnd bool) int {
	index := strings.Index(content, end)
	if index == -1 {
		out.WriteString(content)
		return len(content)
	}

	copied := index
	if includeEnd {
		copied += len(end)
	}

	out.WriteString(content[:copied])

	return copied
}
