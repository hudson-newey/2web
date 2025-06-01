package pageBuilder

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	"strings"
)

type nodeType = int

const (
	htmlNode nodeType = iota
	jsNode
	cssNode
	codeNode
)

func BuildPage(content string) page.Page {
	nonHtmlStartTags := []string{"<script", "<style", "<code"}
	nonHtmlEndTags := []string{"</script>", "</style>", "</code>"}

	currentNodeType := htmlNode
	bufferedContent := ""

	pageModel := page.Page{}
	pageModel.Html = &html.HTMLFile{}

	for i := range content {
		if currentNodeType != htmlNode {
			for _, endTag := range nonHtmlEndTags {
				if i-len(endTag) < 0 {
					continue
				}

				if content[i-len(endTag):i] == endTag {
					if bufferedContent != "" {
						if currentNodeType == jsNode {
							newJsNode := javascript.JSFile{}
							newJsNode.AddContent(bufferedContent)

							pageModel.AddScript(&newJsNode)
						} else if currentNodeType == cssNode {
							newCssNode := css.CSSFile{}
							newCssNode.AddContent(bufferedContent)

							pageModel.AddStyle(&newCssNode)
						} else if currentNodeType == codeNode {
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
					case "<style":
						currentNodeType = cssNode
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
		case jsNode, cssNode, codeNode:
			bufferedContent += string(content[i])
		}
	}

	return pageModel
}
