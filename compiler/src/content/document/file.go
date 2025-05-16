package document

import (
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/html"
	"hudson-newey/2web/src/content/javascript"
)

type Page struct {
	Html       []html.HTMLFile
	JavaScript []javascript.JSFile
	Css        []css.CSSFile
}
