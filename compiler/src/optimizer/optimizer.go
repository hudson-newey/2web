package optimizer

import (
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/optimizer/minify"
)

func OptimizePage(pageModel *page.Page) {
	pageModel.Html.Content = minify.MinifyHtml(pageModel.Html.Content)

	for _, cssModel := range pageModel.Css {
		cssModel.Content = minify.MinifyCss(cssModel.Content)
	}

	for _, jsModel := range pageModel.JavaScript {
		jsModel.Content = minify.MinifyJs(jsModel.Content)
	}
}
