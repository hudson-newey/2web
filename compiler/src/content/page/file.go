package page

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"os"
	"path/filepath"
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

func (model *Page) Write(fileName string) {
	writeFile(fileName+".html", model.html.Content)

	for _, file := range model.css {
		writeFile(file.Content, file.FileName())
	}

	for _, file := range model.javaScript {
		writeFile(file.Content, file.FileName())
	}
}

func (model *Page) AsComponent() {
}

func writeFile(content string, outputPath string) {
	if *cli.GetArgs().ToStdout {
		fmt.Println(content)
	} else {
		os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		os.WriteFile(outputPath, []byte(content), 0644)
	}
}
