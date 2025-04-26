package optimizer

import "strings"

func minifyHtml(content string) string {
	// if the content is less than 2, then the document cannot have a doctype
	// and the content is really small
	if len(content) < 2 {
		return content
	}

	// If the content contains pre elements, we cannot minify the content
	// otherwise the pre-formatted content will be destroyed.
	// This check is actually overly eager, and can trigger even when there are
	// no <pre> elements.
	// But, I would prefer to skip some optimization over breaking the website.
	//
	// TODO: We should add support for minifying with <pre> elements
	hasPre := strings.Contains(content, "<pre")
	if hasPre {
		return content
	}

	// TODO: this will break elements preserving content whitespace
	minifiedContent := strings.ReplaceAll(content, "\n", "")

	return minifiedContent
}
