package page

import (
	"fmt"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
	"strings"
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

	// If we have a full document, we want to inject the script into the <head>
	// element.
	// However, if there is no <head> element (e.g. for a component that is not a
	// page), we want to append it to the end.
	injectLevel := document.Trailing
	if strings.Contains(model.Html.Content, "<head>") {
		injectLevel = document.Head
	}

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
		// type="module" is automatically deferred, but we also add the "async"
		// attribute to the script tag so that the reactive javascript only starts
		// loading when the main thread is free.
		// This will allow the user to start reading/navigating on the page before
		// the reactive content has loaded.
		injectedContent = fmt.Sprintf(`<script async type="module" src="%s"></script>%s`, jsFile.FileName(), "\n")
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

	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, injectLevel)
}

func (model *Page) AddStyle(cssFile *css.CSSFile) {
	model.Css = append(model.Css, cssFile)

	// If we have a full document, we want to inject the styles into the <head>
	// element.
	// However, if there is no <head> element (e.g. for a component that is not a
	// page), we want to append it to the end.
	injectLevel := document.Trailing
	if strings.Contains(model.Html.Content, "<head>") {
		injectLevel = document.Head
	}

	// Adds a "<link>" tag to the html document to load the css file
	injectedContent := fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, cssFile.FileName())
	model.Html.Content = document.InjectContent(model.Html.Content, injectedContent, injectLevel)
}
