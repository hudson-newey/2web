package html

import (
	"fmt"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/javascript"
)

type htmlCode = string

type HTMLFile struct {
	Content htmlCode
}

func (model *HTMLFile) AddContent(contentPartial htmlCode) {
	model.Content += contentPartial
}

// Adds a "<link>" tag to the html document to load the css file
func (model *HTMLFile) AttachCss(file css.CSSFile) {
	injectedContent := fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, file.FileName())
	model.Content = document.InjectContent(model.Content, injectedContent, document.Head)
}

// Adds a "<script src=></script>" tag to the html document to load the js file
func (model *HTMLFile) AttachJs(file javascript.JSFile) {
	// If the JavaScript is not lazy loaded, we want to eagerly evaluate it by not
	// using the "async" keyword.
	// Note that this will delay the initial page load times because all of the
	// JavaScript will be blocking.
	injectedContent := ""
	if file.IsLazy() {
		// type="module" is automatically deferred, but we also add the "async"
		// attribute to the script tag so that the reactive javascript only starts
		// loading when the main thread is free.
		// This will allow the user to start reading/navigating on the page before
		// the reactive content has loaded.
		injectedContent = fmt.Sprintf(`<script async type="module" src="%s"></script>`, file.FileName())
	} else {
		// Note that the script is still deferred because it is using type="module"
		// I have purposely done this so that the programmer doesn't have to deal
		// with differing import formats when using eager/lazy scripting.
		// Eager/lazy loading should be an implementation detail transparent to the
		// programmer.
		// This also allows the user to reference DOM elements in their scripts
		// without any iife or any other related hackery.
		injectedContent = fmt.Sprintf(`<script type="module" src="%s"></script>`, file.FileName())
	}

	model.Content = document.InjectContent(model.Content, injectedContent, document.Head)
}
