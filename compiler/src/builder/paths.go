package builder

import (
	twoWeb "hudson-newey/2web/src/content/2web"
	"hudson-newey/2web/src/content/docx"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/odt"
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
	//
	// Note that we preserve the .xhtml suffix for XHTML files, because they
	// can be natively rendered by browsers.
	adjustedFileName := fileName
	if markdown.IsMarkdownFile(fileName) {
		adjustedFileName = strings.TrimSuffix(adjustedFileName, ".md") + ".html"
	} else if twoWeb.IsTwoWebFile(fileName) {
		adjustedFileName = strings.TrimSuffix(adjustedFileName, ".2web") + ".html"
	} else if docx.IsDocxFile(fileName) {
		adjustedFileName = strings.TrimSuffix(adjustedFileName, ".docx") + ".html"
	} else if odt.IsOdtFile(fileName) {
		adjustedFileName = strings.TrimSuffix(adjustedFileName, ".odt") + ".html"
	}

	adjustedFileName = path.Base(adjustedFileName)
	adjustedPath := path.Join(path.Dir(outputPath), adjustedFileName)

	return adjustedPath
}
