package documentErrors

import (
	_ "embed"
)

//go:embed static/errorTemplate.html
var errorTemplate string

func errorHtmlSource() string {
	return errorTemplate
}
