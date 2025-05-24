package css

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hudson-newey/2web/src/cli"
	"strings"
)

type cssCode = string

type CSSFile struct {
	Content          cssCode
	memoisedFileName string
}

func (model *CSSFile) RawContent() string {
	result := model.Content

	result = strings.ReplaceAll(result, "<style>", "")
	result = strings.ReplaceAll(result, "<style compiled>", "")
	result = strings.ReplaceAll(result, "</style>", "")

	return result
}

func (model *CSSFile) AddContent(newContent cssCode) {
	model.Content += newContent
}

// A css file index that can be used for development builds.
// This should not be used during production builds as it may lead to stale data
// being served from a cdn or browser cache.
var cssFileIndex int = 0

func (model *CSSFile) OutputPath() string {
	outPath := *cli.GetArgs().OutputPath
	return fmt.Sprintf("%s/%s", outPath, model.FileName())
}

// The file hash can be used to uniquely import this file.
// Warning: If the content of this file is changed, the file name will also
// change.
func (model *CSSFile) FileName() string {
	if model.memoisedFileName != "" {
		return model.memoisedFileName
	}

	// If we are in development mode, we want to optimize for build times.
	// We therefore do not compute the md5 hash of the file (which is
	// computationally expensive), and instead just use an incrementing number.
	// This is less efficient for the CDN's and browser cache, but provides a
	// quicker development environment.
	//
	// TODO: remove this && false and fix browser caching
	if !*cli.GetArgs().IsProd && false {
		cssFileIndex++
		result := fmt.Sprintf("%d.css", cssFileIndex)

		model.memoisedFileName = result
		return result
	}

	hash := md5.Sum([]byte(model.Content))

	// Only return the first 8 characters to prevent overly long file names.
	// Warning: This greatly increases the probability of a hash collision.
	// TODO: Keep track of all the used hashes
	fileHash := hex.EncodeToString(hash[:8])

	result := fmt.Sprintf("%s.css", fileHash)

	model.memoisedFileName = result
	return result
}
