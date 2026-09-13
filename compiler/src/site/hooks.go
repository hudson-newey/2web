package site

import (
	"hudson-newey/2web/src/cli"
	"sort"
	debugjson "hudson-newey/2web/src/site/debug.json"
	robotstxt "hudson-newey/2web/src/site/robots.txt"
	sitemapxml "hudson-newey/2web/src/site/sitemap.xml"
)

func AfterAll() {
	paths := GetSitePaths()

	// Pages are registered in compilation completion order, which is
	// nondeterministic when pages are compiled in parallel. Sorting the paths
	// ensures that generated site assets (e.g. the sitemap) are byte-for-byte
	// identical between serial and parallel builds.
	sort.Strings(paths)

	containsSitemap := pathsContain(paths, "sitemap.xml")
	if !containsSitemap {
		sitemapxml.GenerateSitemap(paths)
	}

	containsRobotsTxt := pathsContain(paths, "robots.txt")
	if !containsRobotsTxt {
		robotstxt.GenerateRobotsTxt()
	}

	// We don't want to publish debug info to production as it might contain
	// sensitive information about the source code.
	if !cli.GetArgs().IsProd {
		debugjson.GenerateDebugJson()
	}
}

func BeforeEach(outputPath string) {
	registerSitePage(outputPath)
}

func pathsContain(paths []string, target string) bool {
	for _, path := range paths {
		if path == target {
			return true
		}
	}

	return false
}
