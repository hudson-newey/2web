package optimizer

import (
	"hudson-newey/2web/src/content/page"
)

func injectRuntimeOptimizations(page *page.Page) {
	injectContainerOptimizations(page)
}
