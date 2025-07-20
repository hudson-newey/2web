package runtimeOptimizer

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/constants"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	"strings"
)

func InjectRuntimeOptimizations(page page.Page) page.Page {
	args := cli.GetArgs()

	// TODO: Remove the hack: We have a space before the JS element namespace so
	// that only elements match the replacement target, and not javascript code
	// that uses the selector.
	elementSelector := " " + javascript.JsElementNamespace
	reactiveElementCount := strings.Count(page.Html.Content, elementSelector)
	if reactiveElementCount == 0 {
		return page
	}

	// Note that because container optimizations require shipping runtime css,
	// it increases the bundle size (by a very small amount).
	//
	// However, this very small increase in bundle size is justified because it
	// improves the responsiveness of the page.
	//
	// During development builds, we give this container a human-readable name
	// "isolated-container".
	// However, if we are in a production build, we shorten the css class name.
	isolatedContainerClassName := fmt.Sprintf("%sisolated-container", constants.CompilerNamespace)
	if *args.IsProd {
		isolatedContainerClassName = fmt.Sprintf("%so0", constants.CompilerNamespace)
	}

	runtimeStyles := fmt.Sprintf(`
		.%s{
			contain: content;
		}
	`, isolatedContainerClassName)

	newStyleSheet := css.CSSFile{}
	newStyleSheet.AddContent(runtimeStyles)

	page.AddStyle(&newStyleSheet)

	// Because the 2Web compiler knows exactly what is going to change (down to
	// specific characters in a sentence), and when it going to change.
	// We can put these reactive portions of the page into their own rendering
	// context, so that if something updates, on the page, we can optimize
	// rendering down to the minimum number of pixels that need to be re-painted.
	//
	// Every time we find a reactive element that's traced by the compiler, add
	// the "isolated-container" style to it.
	page.Html.Content = strings.ReplaceAll(
		page.Html.Content,
		elementSelector,
		fmt.Sprintf(
			" class=\"%s\"%s",
			isolatedContainerClassName,
			elementSelector,
		),
	)

	return page
}
