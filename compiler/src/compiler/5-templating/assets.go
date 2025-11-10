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
	if file, err := os.ReadFile(routeCssFile); err == nil {
		cssFile := css.FromContent(string(file))
		page.AddStyle(cssFile)
	}

	// You can only have a single script file per route, so we check for both
	// __script.js and __script.ts, and add them if they exist.
	// TODO: This should probably expand when we add support for other scripting
	// languages.

	routeTsFile := fmt.Sprintf("%s/__script.ts", directory)
	if file, err := os.ReadFile(routeTsFile); err == nil {
		tsFile := javascript.FromContent(string(file))
		page.AddScript(tsFile)
	} else {
		routeJsFile := fmt.Sprintf("%s/__script.js", directory)
		if file, err := os.ReadFile(routeJsFile); err == nil {
			jsFile := javascript.FromContent(string(file))
			page.AddScript(jsFile)
		}
	}
}
