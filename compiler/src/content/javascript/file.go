package javascript

type javascriptCode = string

type JSFile struct {
	content javascriptCode

	Eager bool
}

func (model *JSFile) AddContent(partialContent string) {
	model.content += partialContent
}

// Whether to eagerly execute the JavaScript code. Eagerly executed code is
// typically reserved for JavaScript code that effects the layout of the page.
//
// E.g. Some JavaScript that needs to run to set a property on an element that
// is not exposed through an attribute.
func (model *JSFile) IsEager() bool {
	return false
}
