package builder

import (
	"hudson-newey/2web/src/content/markdown"
	"os"
	"path"
	"strings"
)

func outputFileName(inputPath string, outputPath string, fileName string) string {
	isInDir := strings.HasSuffix(inputPath, string(os.PathSeparator))
	isOutDir := strings.HasSuffix(outputPath, string(os.PathSeparator))

	// If we are passed in a file, but a directory as the output, append the
	// file name to the output directory
	if !isInDir && isOutDir {
		fileName = outputPath + fileName
	}

	// If we are compiling a markdown file, we want to replace the .md suffix
	// with .html
	// This is because we compile markdown to html files.
	adjustedFileName := fileName
	if markdown.IsMarkdownFile(fileName) {
		adjustedFileName = strings.TrimSuffix(adjustedFileName, ".md") + ".html"
	}

	adjustedFileName = path.Base(adjustedFileName)
	adjustedPath := path.Join(path.Dir(outputPath), adjustedFileName)

	return adjustedPath
}
