package javascript

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hudson-newey/2web/src/cli"
)

type javascriptCode = string

type JSFile struct {
	Content javascriptCode
}

func (model *JSFile) AddContent(partialContent string) {
	model.Content += partialContent
}

// Whether to eagerly execute the JavaScript code. Eagerly executed code is
// typically reserved for JavaScript code that effects the layout of the page.
//
// E.g. Some JavaScript that needs to run to set a property on an element that
// is not exposed through an attribute.
func (model *JSFile) IsLazy() bool {
	return true
}

// A js file index that can be used for development builds.
// This should not be used during production builds as it may lead to stale data
// being served from a cdn or browser cache.
var jsFileIndex int = 0

// The file hash can be used to uniquely import this file.
// Warning: If the content of this file is changed, the file name will also
// change.
func (model *JSFile) FileName() string {
	// If we are in development mode, we want to optimize for build times.
	// We therefore do not compute the md5 hash of the file (which is
	// computationally expensive), and instead just use an incrementing number.
	// This is less efficient for the CDN's and browser cache, but provides a
	// quicker development environment.
	if !*cli.GetArgs().IsProd {
		jsFileIndex++
		return fmt.Sprint(jsFileIndex)
	}

	hash := md5.Sum([]byte(model.Content))

	// Only return the first 8 characters to prevent overly long file names.
	// Warning: This greatly increases the probability of a hash collision.
	// TODO: Keep track of all the used hashes
	return hex.EncodeToString(hash[:8])
}
