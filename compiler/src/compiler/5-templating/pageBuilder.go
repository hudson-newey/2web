package templating

import (
	"fmt"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	"os"
	"path/filepath"
	"strings"
)

type languageState = int

const (
	htmlNode languageState = iota
	jsNode
	codeNode
)

func BuildPage(filePath string, content string) page.Page {
	nonHtmlStartTags := []string{"<script", "<code"}
	nonHtmlEndTags := []string{"</script>", "</code>"}

	currentNodeType := htmlNode
	bufferedContent := ""

	pageModel := page.Page{InputPath: filePath}
	pageModel.Html = &html.HTMLFile{}

	for i := range content {
		if currentNodeType != htmlNode {
			for _, endTag := range nonHtmlEndTags {
				if i-len(endTag) < 0 {
					continue
				}

				if content[i-len(endTag):i] == endTag {
					if bufferedContent != "" {
						if currentNodeType == codeNode {
							contentToPrepend := strings.TrimPrefix(bufferedContent, "<code>")

							escapedContent := html.EscapeHtml(contentToPrepend)
							escapedContent = "<code>" + escapedContent

							pageModel.Html.AddContent(escapedContent)
						}
					}

					currentNodeType = htmlNode
					bufferedContent = ""
					break
				}
			}
		}

		// We do not allow transitioning to other tag types if we are in a code node
		// so that you can write script and style tags inside of the code block.
		// We escape all of the tags in a later stage.
		if currentNodeType != codeNode {
			for _, startTag := range nonHtmlStartTags {
				if i+len(startTag) > len(content) {
					continue
				}

				if content[i:i+len(startTag)] == startTag {
					switch startTag {
					case "<script":
						currentNodeType = jsNode
					case "<code":
						currentNodeType = codeNode
					}
				}
			}
		}

		// We purposely ignore the compiledJsNode case because we don't want it
		// included in the production build.
		switch currentNodeType {
		case htmlNode:
			pageModel.Html.AddContent(string(content[i]))
		case jsNode, codeNode:
			bufferedContent += string(content[i])
		}
	}

	addRouteAssets(&pageModel)

	return pageModel
}

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
		jsFile := javascript.FromFilePath(routeTsFile)
		page.AddScript(jsFile)
	}
}
