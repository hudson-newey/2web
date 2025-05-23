package svg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hudson-newey/2web/src/cli"
)

type svgCode = string

type SvgFile struct {
	Content svgCode
}

func (model *SvgFile) AddContent(newContent string) {
	model.Content += newContent
}

// A css file index that can be used for development builds.
// This should not be used during production builds as it may lead to stale data
// being served from a cdn or browser cache.
var svgFileIndex int = 0

// The file hash can be used to uniquely import this file.
// Warning: If the content of this file is changed, the file name will also
// change.
func (model *SvgFile) FileName() string {
	// If we are in development mode, we want to optimize for build times.
	// We therefore do not compute the md5 hash of the file (which is
	// computationally expensive), and instead just use an incrementing number.
	// This is less efficient for the CDN's and browser cache, but provides a
	// quicker development environment.
	if !*cli.GetArgs().IsProd {
		svgFileIndex++
		return fmt.Sprintf("%d.css", svgFileIndex)
	}

	hash := md5.Sum([]byte(model.Content))

	// Only return the first 8 characters to prevent overly long file names.
	// Warning: This greatly increases the probability of a hash collision.
	// TODO: Keep track of all the used hashes
	return hex.EncodeToString(hash[:8]) + ".css"
}
