package javascript

type javascriptCode = string

type JSFile struct {
	Content javascriptCode

	// Whether to eagerly execute the JavaScript code. Eagerly executed code is
	// typically reserved for JavaScript code that effects the layout of the page.
	//
	// E.g. Some JavaScript that needs to run to set a property on an element that
	// is not exposed through an attribute.
	Eager bool
}
