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
	// We make all of the css rules important to ensure they override any user
	// styles and so that the user cannot accidentally break the error page
	// styles by indirectly adding styles to elements like a h2 or p tag.
	templateCssWithOverrides := strings.ReplaceAll(
		errorTemplateCSS,
		";",
		" !important;",
	)

	// To include the CSS in the HTML template, I hacked together a solution
	// where I replace a substring placeholder with the actual CSS content.
	//
	// TODO: Fix this properly in the future with a correct templating approach.
	replacementTarget := "/* __2web:compiler-internal error-template-styles */"
	finalTemplate := strings.Replace(
		errorTemplate,
		replacementTarget,
		templateCssWithOverrides,
		1,
	)

	return finalTemplate
}
