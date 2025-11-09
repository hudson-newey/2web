package templating

import "hudson-newey/2web/src/content/page"

// BuildComponentPage is a function injected by the builder package to compile
// a component file path into a page.Page. It is declared here to avoid a
// direct import on the builder package which would cause an import cycle.
//
// The default implementation returns (empty, false) indicating the build could
// not be performed. The builder package should assign this to its BuildToPage
// function during init().
var BuildComponentPage = func(inputPath string, isFullPage bool) (page.Page, bool) {
	return page.NewPage(), false
}
