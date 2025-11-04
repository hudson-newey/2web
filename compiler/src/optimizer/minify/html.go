package minify

import (
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

// I use a minification library because third party minification libraries will
// always be able to get a better minification result than a custom
// implementation.
// A custom implementation would save further file content, but a custom
// minifier can be bootstrapped on top of the library.

// Minifies HTML content, including all inline styles, scripts, and svg's
func MinifyHtml(content string) string {
	// if the content is less than 2, then the document cannot have a doctype
	// and the content is really small
	if len(content) < 2 {
		return content
	}

	// keep document tags such as <head> because they help the browser prefetcher
	// perform dns queries
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepWhitespace:   false,
	})
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	// For some reason, the minification library that I am using doesn't like
	// <code> blocks inside of html.
	// It incorrectly tries to format <script> tags inside of <code> blocks as if
	// it was real JavaScript code.
	//
	// However, code inside of <code> example blocks are user provided and we
	// should not be minifying them.
	// TODO: find out if this is my fault for including a plugin that minifies
	// <code> blocks.
	//
	// Additionally, the minifier will throw errors if the code inside the fake
	// <script> blocks contain invalid JavaScript.
	// This is slightly annoying because I believe that code examples can contain
	// invalid JavaScript code.
	if !strings.Contains(content, "</code>") {
		m.AddFunc("application/javascript", js.Minify)
	}

	minifiedContent, err := m.String("text/html", content)
	if err != nil {
		panic(err)
	}

	// The "minify" library can get 99% of the way to full minification, however
	// there are some extreme minification techniques that we can perform to get
	// assets even smaller.
	result := customMinifier(minifiedContent)
	return result
}

func customMinifier(content string) string {
	// If the content contains pre elements, we cannot minify the content
	// otherwise the pre-formatted content will be destroyed.
	// This check is actually overly eager, and can trigger even when there are
	// no <pre> elements.
	// But, I would prefer to skip some optimization over breaking the website.
	//
	// TODO: We should add support for minifying with <pre> elements
	hasPre := strings.Contains(content, "<pre")
	if hasPre {
		return content
	}

	// TODO: this will break elements preserving content whitespace
	minifiedContent := strings.ReplaceAll(content, "\n", "")

	return minifiedContent
}
