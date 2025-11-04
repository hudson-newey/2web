package builder

import (
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/page"
)

func isBinaryData(data []byte) bool {
	// A simple heuristic to check for binary data:
	// If the data contains a null byte, we consider it binary.
	for _, b := range data {
		if b == 0 {
			return true
		}
	}

	return false
}

func buildFromBinary(inputPath string, data []byte) (page.Page, bool) {
	compiledPage := page.Page{}
	compiledPage.Html = &html.HTMLFile{}
	compiledPage.Html.AddContent("[Binary content cannot be displayed]")

	isErrorFree := true

	return compiledPage, isErrorFree
}
