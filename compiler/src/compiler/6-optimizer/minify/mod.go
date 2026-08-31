package minify

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/page"
)

func MinifyPage(pageModel *page.Page) {
	args := cli.GetArgs()

	// If we are not doing a production build, don't spend cycles optimizing
	// code that doesn't need to be highly optimized.
	if !args.IsProd {
		return
	}

	if args.WithFormatting {
		cli.PrintWarning("Ignoring '--format' because '--production' was specified")
	}

	pageModel.Html.Content = minifyHtml(pageModel.Html.Content)

	for _, cssModel := range pageModel.Css {
		cssModel.Content = minifyCss(cssModel.Content)
	}

	for _, jsModel := range pageModel.JavaScript {
		jsModel.Content = minifyJs(jsModel.Content)
	}
}
