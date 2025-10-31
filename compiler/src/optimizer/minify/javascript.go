package minify

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

func MinifyJs(content string) string {
	m := minify.New()

	// Check if the "uglifyjs" command-line tool is available
	// and use it for minification if so.
	// Fallback to the built-in minifier otherwise.
	if _, err := exec.LookPath("uglifyjs"); err == nil {
		m.AddCmdRegexp(
			regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"),
			exec.Command("uglifyjs"),
		)
	} else {
		fmt.Println("Warning: 'uglifyjs' not found in PATH. Using the built-in minifier may cause larger file sizes.")
		m.AddFunc("application/javascript", js.Minify)
	}

	minifiedContent, err := m.String("application/javascript", content)
	if err != nil {
		panic(err)
	}

	return minifiedContent
}
