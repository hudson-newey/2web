package runtimeOptimizer

import (
	"hudson-newey/2web/src/content/page"
)

func InjectRuntimeOptimizations(page *page.Page) {
	injectContainerOptimizations(page)
}
