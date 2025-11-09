package documentErrors

import (
	_ "embed"
	"strings"
)

//go:embed static/errorTemplate.html
var errorTemplate string

//go:embed static/errorTemplate.css
var errorTemplateCSS string

func errorHtmlSource() string {
	// To include the CSS in the HTML template, I hacked together a solution
	// where I replace a substring placeholder with the actual CSS content.
	//
	// TODO: Fix this properly in the future with a correct templating approach.
	replacementTarget := "/* __2web:compiler-internal error-template-styles */"
	finalTemplate := strings.Replace(
		errorTemplate,
		replacementTarget,
		errorTemplateCSS,
		1,
	)

	return finalTemplate
}
