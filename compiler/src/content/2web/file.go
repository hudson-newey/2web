package twoWeb

import "hudson-newey/2web/src/content/html"

// At the moment, the .2web file format is processed the exact same as a
// .html file.
// However, I have used a different type here to allow for future
// differentiation if needed.
type TwoWebFile = html.HTMLFile
