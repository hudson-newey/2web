package minify

import (
	"github.com/hudson-newey/2web/_shared/logger"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

func minifyCss(content string) string {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	minifiedContent, err := m.String("text/css", content)
	if err != nil {
		// A minification failure must not kill the entire build.
		logger.PrintWarning("failed to minify css, shipping unminified content: " + err.Error())
		return content
	}

	return minifiedContent
}
