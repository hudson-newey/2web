package optimizer

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler/6-optimizer/minify"
	"hudson-newey/2web/src/content/page"
)

func OptimizePage(pageModel *page.Page) {
	args := cli.GetArgs()

	// We perform runtime optimizations here so that they are only applied once.
	// If we instead applied them inside of the BuildToPage function, then the
	// optimizations would be applied every time a component is added.
	// Most optimizations can be applied at the very end of the page build, so
	// this is the best place to do it.
	if !(args.NoRuntimeOptimizations) {
		injectRuntimeOptimizations(pageModel)
	}

	pageModel.Html.Content = minify.MinifyHtml(pageModel.Html.Content)

	for _, cssModel := range pageModel.Css {
		cssModel.Content = minify.MinifyCss(cssModel.Content)
	}

	for _, jsModel := range pageModel.JavaScript {
		jsModel.Content = minify.MinifyJs(jsModel.Content)
	}
}
