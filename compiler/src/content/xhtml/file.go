package xhtml

import (
	"hudson-newey/2web/src/content/html"
)

// TODO: XHTML files are technically different from HTML files, so we probably
// do want to add a specialized XHTMLFile type in the future.
// However, for now, this is sufficient, because it means that we can share the
// same document partial and templating logic between HTML and XHTML files.
type XHTMLFile = html.HTMLFile
