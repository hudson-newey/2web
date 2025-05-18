package document

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

type Page struct {
	html       *html.HTMLFile
	javaScript []*javascript.JSFile
	css        []*css.CSSFile
}

func (model *Page) SetContent(htmlFile *html.HTMLFile) {
	model.html = htmlFile
}

func (model *Page) AddScript(jsFile *javascript.JSFile) {
	model.javaScript = append(model.javaScript, jsFile)
}

func (model *Page) AddStyle(cssFile *css.CSSFile) {
	model.css = append(model.css, cssFile)
}

func (model *Page) Write(path string) {
}
