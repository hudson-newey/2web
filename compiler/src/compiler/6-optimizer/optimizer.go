package optimizer

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/compiler/6-optimizer/assetOptimization"
	"hudson-newey/2web/src/compiler/6-optimizer/minify"
	"hudson-newey/2web/src/compiler/6-optimizer/runtimeOptimization"
	"hudson-newey/2web/src/content/page"
)

func OptimizePage(pageModel *page.Page) {
	// If we are not doing a production build, don't spend cycles optimizing
	// code that doesn't need to be highly optimized.
	args := cli.GetArgs()
	if !args.IsProd {
		return
	}
	if args.WithFormatting {
		cli.PrintWarning("Ignoring '--format' because optimizations are enabled")
	}

	runtimeOptimization.InjectRuntimeOptimizations(pageModel)
	assetOptimization.OptimizeAssets(pageModel)
	// Minify after runtime optimizations so that the additional classes +
	// markup can be minified.
	minify.MinifyPage(pageModel)
}
