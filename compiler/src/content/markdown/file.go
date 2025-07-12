package markdown

import (
	"hudson-newey/2web/src/content/html"
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
	markdownHtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type markdownCode = string

type MarkdownFile struct {
	Content markdownCode
}

func (model *MarkdownFile) ToHtml() html.HTMLFile {
	file := html.HTMLFile{Content: ""}

	extensions := parser.CommonExtensions | parser.Titleblock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(model.Content))

	// create HTML renderer with extensions
	htmlFlags := markdownHtml.CommonFlags | markdownHtml.HrefTargetBlank
	opts := markdownHtml.RendererOptions{Flags: htmlFlags}
	renderer := markdownHtml.NewRenderer(opts)

	htmlBytes := markdown.Render(doc, renderer)
	file.AddContent(string(htmlBytes))

	return file
}

func (model *MarkdownFile) Reader() io.Reader {
	return strings.NewReader(model.Content)
}
