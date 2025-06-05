package builder

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/svg"
	"os"
)

func indexPages(inputPath string) []string {
	inputFile, err := os.Stat(inputPath)
	if err != nil {
		panic(err)
	}

	if inputFile.IsDir() {
		totalFiles := []string{}
		currentDirFiles, err := os.ReadDir(inputPath)
		if err != nil {
			panic(err)
		}

		for _, file := range currentDirFiles {
			pages := indexPages(inputPath + "/" + file.Name())

			for _, page := range pages {
				// TODO: Don't include assets in page indexing. They should instead be
				// pulled out of the page source so that they can be efficiently tree
				// shaken.
				if html.IsPage(page) || markdown.IsMarkdownFile(page) || css.IsCssFile(page) || javascript.IsJsFile(page) || svg.IsSvgFile(page) {
					totalFiles = append(totalFiles, page)
				}
			}
		}

		return totalFiles
	}

	return []string{inputPath}
}
