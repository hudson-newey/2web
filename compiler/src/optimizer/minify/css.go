package minify

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

func MinifyCss(content string) string {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	minifiedContent, err := m.String("text/css", content)
	if err != nil {
		panic(err)
	}

	return minifiedContent
}
