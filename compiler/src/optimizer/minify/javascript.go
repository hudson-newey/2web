package minify

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

func MinifyJs(content string) string {
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)

	minifiedContent, err := m.String("text/css", content)
	if err != nil {
		panic(err)
	}

	return minifiedContent
}
