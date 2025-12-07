package builder

import (
	"hudson-newey/2web/src/content/assets"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/svg"
	"os"
	"strings"
)

func indexPages(inputPath string) []string {
	inputFile, err := os.Stat(inputPath)
	if err != nil {
		panic(err)
	}

	if inputFile.IsDir() {
		dirInputPath := inputPath
		if !strings.HasSuffix(inputPath, string(os.PathSeparator)) {
			dirInputPath += string(os.PathSeparator)
		}

		currentDirFiles, err := os.ReadDir(inputPath)
		if err != nil {
			panic(err)
		}

		// Pre-allocate slice with estimated capacity to reduce reallocations.
		// Multiplier of 2 accounts for recursive directory traversal and assumes
		// roughly equal numbers of directories and files at each level.
		estimatedCapacity := len(currentDirFiles) * 2
		totalFiles := make([]string, 0, estimatedCapacity)

		for _, file := range currentDirFiles {
			pages := indexPages(dirInputPath + file.Name())

			for _, page := range pages {
				// TODO: Instead of copying all files, we should only copy files that
				// are referenced/linked directly/indirectly from markdown files.
				shouldPreserve := assets.IsMarkupFile(page) ||
					css.IsCssFile(page) ||
					javascript.IsJsFile(page) ||
					svg.IsSvgFile(page)

				// TODO: Don't include assets in page indexing. They should instead be
				// pulled out of the page source so that they can be efficiently tree
				// shaken.
				if shouldPreserve && !assets.IsComponent(page) {
					totalFiles = append(totalFiles, page)
				}
			}
		}

		// Filter out all paths that are not markup files.
		// This means that any non-markup files will be tree-shaken if they are not
		// used by any markup files.
		// Pre-allocate with the current length as upper bound
		filteredFiles := make([]string, 0, len(totalFiles))
		for _, file := range totalFiles {
			if assets.IsMarkupFile(file) {
				filteredFiles = append(filteredFiles, file)
			}
		}

		return filteredFiles
	}

	return []string{inputPath}
}
