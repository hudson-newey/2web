package builder

import (
	"hudson-newey/2web/src/content/page"
	"unicode/utf8"
)

func BuildToPage(inputPath string, isFullPage bool) (page.Page, bool) {
	data := getContent(inputPath)

	// This does limit us to UTF-8 which will exclude a lot of legacy encodings
	// and character sets from different countries.
	// However, I have decided to default to UTF-8 for the initial version.
	//
	// TODO: We should add support for other encodings in the future.
	if utf8.Valid(*data) {
		return buildFromString(inputPath, string(*data), isFullPage)
	} else {
		return buildFromBinary(inputPath, data, isFullPage)
	}
}
