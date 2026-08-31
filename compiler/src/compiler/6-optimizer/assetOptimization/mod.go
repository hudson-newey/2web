package assetOptimization

import "hudson-newey/2web/src/content/page"

func OptimizeAssets(page *page.Page) {
	compressImages(page)
}
