package templating

import (
	"fmt"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	"os"
	"path/filepath"
)

// TODO: I should find a better way to do this instead of hardcoding the file
// names.
func addRouteAssets(page *page.Page) {
	// If there is a __style.css sidecar file, we want to include it in the page.
	// To do this, we remove the complete file name from the file path so that
	// only the directory remains, and then we append the __style.css file name.
	directory := filepath.Dir(page.InputPath)
	routeCssFile := fmt.Sprintf("%s/__style.css", directory)

	if _, err := os.Stat(routeCssFile); err == nil {
		cssFile := css.FromFilePath(routeCssFile)
		page.AddStyle(cssFile)
	}

	routeJsFile := fmt.Sprintf("%s/__script.js", directory)
	if _, err := os.Stat(routeJsFile); err == nil {
		jsFile := javascript.FromFilePath(routeJsFile)
		page.AddScript(jsFile)
	}

	routeTsFile := fmt.Sprintf("%s/__script.ts", directory)
	if _, err := os.Stat(routeTsFile); err == nil {
		tsFile := javascript.FromFilePath(routeTsFile)
		page.AddScript(tsFile)
	}
}
