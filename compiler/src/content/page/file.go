package page

import (
	"fmt"
	"hudson-newey/2web/src/content"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	twoscript "hudson-newey/2web/src/content/twoScript"
)

func NewPage() Page {
	return Page{
		Html:       &html.HTMLFile{},
		TwoScript:  []*twoscript.TwoScriptFile{},
		JavaScript: []*javascript.JSFile{},
		Css:        []*css.CSSFile{},
		Assets:     []*content.BinaryFile{},
	}
}

type Page struct {
	InputPath  string
	Html       *html.HTMLFile
	TwoScript  []*twoscript.TwoScriptFile
	JavaScript []*javascript.JSFile
	Css        []*css.CSSFile
	Assets     []*content.BinaryFile
}

func (model *Page) SetContent(htmlFile *html.HTMLFile) {
	model.Html = htmlFile
}

func (model *Page) AddTwoScript(tsFile *twoscript.TwoScriptFile) {
	model.TwoScript = append(model.TwoScript, tsFile)
}

func (model *Page) AddScript(jsFile *javascript.JSFile) {
	model.JavaScript = append(model.JavaScript, jsFile)

	// Adds a "<script src=></script>" tag to the html document to load the js file
	// If the JavaScript is not lazy loaded, we want to eagerly evaluate it by not
	// using the "async" keyword.
	// Note that this will delay the initial page load times because all of the
	// JavaScript will be blocking.
	//
	// Note that we include new line character after every import to make them
	// easier to read in a development environment.
	// These new line characters will be removed when building for production.
	injectedContent := ""
	if jsFile.IsLazy() {
		// type="module" is automatically deferred.
		injectedContent = fmt.Sprintf(`<script type="module" src="%s"></script>%s`, jsFile.FileName(), "\n")
	} else {
		// Note that the script is still deferred because it is using type="module"
		// I have purposely done this so that the programmer doesn't have to deal
		// with differing import formats when using eager/lazy scripting.
		// Eager/lazy loading should be an implementation detail transparent to the
		// programmer.
		// This also allows the user to reference DOM elements in their scripts
		// without any iife or any other related hackery.
		injectedContent = fmt.Sprintf(`<script type="module" src="%s"></script>%s`, jsFile.FileName(), "\n")
	}

	// While normal (non-module type) scripts block DOM AST construction, because
	// 2Web uses async + module type scripts, execution is deferred until the DOM
	// has been constructed.
	// Meaning that it is actually more beneficial to inject scripts into the top
	// of the head, so that they can be discovered and start fetching sooner.
	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, document.HeadTop)
}

func (model *Page) AddStyle(cssFile *css.CSSFile) {
	model.Css = append(model.Css, cssFile)

	// Adds a "<link>" tag to the html document to load the css file
	injectedContent := fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" />\n", cssFile.FileName())

	// Always inject css styles into the top of the head element so that they can
	// be discovered as soon as possible, to begin parsing
	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, document.HeadTop)
}

// Adds an asset (binary file) to the page model.
// This is an escape hatch for adding binary or unrecognized formats to the
// output assets, and a smell that the compiler cannot correctly handle the file
// type and perform necessary optimizations/transpilation/etc.
func (model *Page) AddAsset(binaryFile *content.BinaryFile) {
	model.Assets = append(model.Assets, binaryFile)
}

func (model *Page) Format() {
	if model.Html != nil {
		model.Html.Format()
	}

	for _, jsFile := range model.JavaScript {
		jsFile.Format()
	}

	for _, cssFile := range model.Css {
		cssFile.Format()
	}
}
