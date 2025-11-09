package builder

import templating "hudson-newey/2web/src/compiler/5-templating"

// This hack exists to get around cyclical imports because the components are
// built using the same builder that is trying to build the main page, and pages
// can import components, causing an import cycle.
// TODO: Replace this with a better solution.
func init() {
	// Wire the component build function to avoid an import cycle. This lets the
	// templating package request component builds without importing builder.
	templating.BuildComponentPage = BuildToPage
}
