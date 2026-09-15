package preprocessor

import (
	"regexp"
	"strings"
)

// Components can declare slots that the consumer of the component fills:
//
//	<!-- component -->
//	<slot name="header"></slot>
//	<p><slot>fallback</slot></p>
//
//	<!-- consumer -->
//	<Counter>
//		<div slot="header">Hello</div>
//		<p>default content</p>
//	</Counter>
//
// The consumer's children that carry a slot="name" attribute replace the
// named slot with that name (the attribute is removed from the rendered
// content), and the remaining children replace the default slot. A slot
// without matching content renders its own fallback content (the markup
// inside the slot tag).

// Matches a slot declaration of a component. e.g. '<slot name="header"></slot>'
var slotPattern = regexp.MustCompile(`(?s)<slot\b([^>]*?)(?:/>|>(.*?)</slot\s*>)`)

// Matches a slot attribute of a consumer's child. e.g. slot="header"
var slotAttributePattern = regexp.MustCompile(`\s*slot="([^"]*)"`)

// Matches the name attribute of a slot declaration.
var slotNamePattern = regexp.MustCompile(`name="([^"]*)"`)

// parseInstance splits a component instance into its attribute string and its
// children (empty for self closing instances).
func parseInstance(selector string) (attributes string, children string) {
	if strings.HasSuffix(selector, "/>") {
		return strings.TrimSuffix(strings.TrimPrefix(selector, "<"), "/>"), ""
	}

	openEnd := strings.Index(selector, ">")

	return selector[1:openEnd], selector[openEnd+1 : strings.LastIndex(selector, "</")]
}

// spliceSlots replaces the slot declarations of a component's content with
// the content the instance passed for them.
func spliceSlots(componentContent string, children string) string {
	named, defaultContent := extractSlotContent(children)

	return slotPattern.ReplaceAllStringFunc(componentContent, func(match string) string {
		groups := slotPattern.FindStringSubmatch(match)

		attributes := groups[1]
		fallback := groups[2]

		nameMatch := slotNamePattern.FindStringSubmatch(attributes)

		if nameMatch == nil {
			// The default slot.
			if strings.TrimSpace(defaultContent) == "" {
				return fallback
			}

			return defaultContent
		}

		if content, ok := named[nameMatch[1]]; ok {
			return content
		}

		return fallback
	})
}

// extractSlotContent splits a component instance's children into the contents
// of the named slots and the default slot content.
//
// The children of an instance are consumer owned markup, so they are scanned
// verbatim: everything that isn't a slot assignment (text, elements without a
// slot attribute) is default slot content, in source order.
func extractSlotContent(children string) (named map[string]string, defaultContent string) {
	named = map[string]string{}

	var defaultParts []string

	i := 0
	pendingStart := 0

	// flushDefault appends the pending content (text and elements that aren't
	// slot assignments) to the default slot content.
	flushDefault := func(end int) {
		if end > pendingStart {
			defaultParts = append(defaultParts, children[pendingStart:end])
		}
	}

	for i < len(children) {
		if children[i] != '<' {
			i++
			continue
		}

		if strings.HasPrefix(children[i:], "<!--") {
			i += 4 + copyUntilMarker(children[i+4:], "-->") + len("-->")
			continue
		}

		tagEnd := findTagEnd(children[i:])
		tag := children[i : i+tagEnd]
		name := tagName(tag)

		slotMatch := slotAttributePattern.FindStringSubmatch(tag)

		isEndTag := strings.HasPrefix(children[i:], "</")
		isSelfClosing := strings.HasSuffix(tag, "/>")

		if slotMatch == nil || isEndTag {
			// Not a slot assignment: it stays in the default slot content.
			i += tagEnd

			if !isEndTag && !isSelfClosing && !voidElements[name] {
				i += matchedElementLength(children[i:], name)
			}

			continue
		}

		flushDefault(i)

		elementLength := tagEnd
		if !isSelfClosing && !voidElements[name] {
			elementLength = matchedElementLength(children[i:], name)
		}

		named[slotMatch[1]] += removeSlotAttribute(children[i:i+elementLength], slotMatch[0])
		i += elementLength
		pendingStart = i
	}

	flushDefault(len(children))

	return named, strings.Join(defaultParts, "")
}

// copyUntilMarker returns the offset of the end marker in the content.
func copyUntilMarker(content string, end string) int {
	index := strings.Index(content, end)
	if index == -1 {
		return len(content)
	}

	return index
}

// isTagNamed returns whether the content starts with a tag that has the given
// name (start or end tag).
func isTagNamed(content string, name string) bool {
	if len(content) < len(name)+2 || content[0] != '<' {
		return false
	}

	rest := content[1:]

	if strings.HasPrefix(rest, "/") {
		rest = rest[1:]
	}

	if len(rest) < len(name) || !strings.EqualFold(rest[:len(name)], name) {
		return false
	}

	return isFollowedByTagClose(rest[len(name):])
}

// isFollowedByTagClose returns whether the content after a tag name is the
// end of the tag name (whitespace, a slash, or ">").
func isFollowedByTagClose(content string) bool {
	if content == "" {
		return true
	}

	switch content[0] {
	case ' ', '\t', '\n', '\r', '>', '/':
		return true
	}

	return false
}

// matchedElementLength returns the length of the element that starts at the
// beginning of the content, matched to the balanced end tag of the same name
// (so that nested elements with the same name don't end the match early).
func matchedElementLength(content string, name string) int {
	depth := 0
	i := 0

	for i < len(content) {
		if !isTagNamed(content[i:], name) {
			i++
			continue
		}

		tagEnd := findTagEnd(content[i:])
		isEndTag := strings.HasPrefix(content[i:], "</")

		if isEndTag {
			i += tagEnd
			depth--

			if depth == 0 {
				return i
			}

			continue
		}

		isSelfClosing := strings.HasSuffix(content[i:i+tagEnd], "/>")

		if !isSelfClosing && !voidElements[name] {
			depth++
		}

		i += tagEnd
	}

	return len(content)
}

// removeSlotAttribute removes a slot attribute from an element's start tag,
// so that the rendered content doesn't carry the assignment anymore.
func removeSlotAttribute(element string, attribute string) string {
	startTagEnd := findTagEnd(element)
	startTag := element[:startTagEnd]

	return strings.Replace(startTag, attribute, "", 1) + element[startTagEnd:]
}
