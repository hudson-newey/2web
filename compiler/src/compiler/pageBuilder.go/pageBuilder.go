package pageBuilder

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
)

type nodeType = int

const (
	htmlNode nodeType = iota
	jsNode
	cssNode
)

func BuildPage(content string) page.Page {
	nonHtmlStartTags := []string{"<script", "<style"}
	nonHtmlEndTags := []string{"</script>", "</style>"}

	codeBlockStart := []string{"<code"}
	codeBlockEnd := []string{"</code>"}

	// When we are inside a code block, we want to emit style and script tag
	// content as if it was typed as text.
	inCodeBlock := false

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
						}
					}

					currentNodeType = htmlNode
					bufferedContent = ""
					break
				}
			}
		}

		for _, startTag := range codeBlockStart {
			if i+len(startTag) > len(content) {
				continue
			}

			if content[i:i+len(startTag)] == startTag {
				inCodeBlock = true
				break
			}
		}

		for _, endTag := range codeBlockEnd {
			if i-len(endTag) < 0 {
				continue
			}

			if content[i-len(endTag):i] == endTag {
				inCodeBlock = false
				break
			}
		}

		if !inCodeBlock {
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
					}
				}
			}
		}

		// We purposely ignore the compiledJsNode case because we don't want it
		// included in the production build.
		switch currentNodeType {
		case htmlNode:
			pageModel.Html.AddContent(string(content[i]))
		case jsNode:
			bufferedContent += string(content[i])
		case cssNode:
			bufferedContent += string(content[i])
		}
	}

	return pageModel
}
