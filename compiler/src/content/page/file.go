package page

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"os"
	"path/filepath"
)

type Page struct {
	Html       *html.HTMLFile
	JavaScript []*javascript.JSFile
	Css        []*css.CSSFile
}

func (model *Page) SetContent(htmlFile *html.HTMLFile) {
	model.Html = htmlFile
}

func (model *Page) AddScript(jsFile *javascript.JSFile) {
	model.JavaScript = append(model.JavaScript, jsFile)

	if jsFile.IsCompilerOnly() {
		return
	}

	// Adds a "<script src=></script>" tag to the html document to load the js file
	// If the JavaScript is not lazy loaded, we want to eagerly evaluate it by not
	// using the "async" keyword.
	// Note that this will delay the initial page load times because all of the
	// JavaScript will be blocking.
	injectedContent := ""
	if jsFile.IsLazy() {
		// type="module" is automatically deferred, but we also add the "async"
		// attribute to the script tag so that the reactive javascript only starts
		// loading when the main thread is free.
		// This will allow the user to start reading/navigating on the page before
		// the reactive content has loaded.
		injectedContent = fmt.Sprintf(`<script async type="module" src="%s"></script>`, jsFile.FileName())
	} else {
		// Note that the script is still deferred because it is using type="module"
		// I have purposely done this so that the programmer doesn't have to deal
		// with differing import formats when using eager/lazy scripting.
		// Eager/lazy loading should be an implementation detail transparent to the
		// programmer.
		// This also allows the user to reference DOM elements in their scripts
		// without any iife or any other related hackery.
		injectedContent = fmt.Sprintf(`<script type="module" src="%s"></script>`, jsFile.FileName())
	}

	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, document.Head)
}

func (model *Page) AddStyle(cssFile *css.CSSFile) {
	model.Css = append(model.Css, cssFile)

	// Adds a "<link>" tag to the html document to load the css file
	injectedContent := fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, cssFile.FileName())
	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, document.Head)
}

func (model *Page) Write(filePath string) {
	writeFile(model.Html.Content, filePath)

	for _, file := range model.Css {
		writeFile(file.RawContent(), file.OutputPath())
	}

	for _, file := range model.JavaScript {
		if file.IsCompilerOnly() {
			continue
		}

		writeFile(file.RawContent(), file.OutputPath())
	}
}

func writeFile(content string, outputPath string) {
	if *cli.GetArgs().ToStdout {
		fmt.Println(content)
	} else {
		os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		os.WriteFile(outputPath, []byte(content), 0644)
	}
}
