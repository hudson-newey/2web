package site

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/models"
	debugjson "hudson-newey/2web/src/site/debug.json"
	robotstxt "hudson-newey/2web/src/site/robots.txt"
	sitemapxml "hudson-newey/2web/src/site/sitemap.xml"
)

func AfterAll() {
	paths := GetSitePaths()

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
	path := models.SitePath{
		Path: outputPath,
	}

	registerSitePage(path)
}

func pathsContain(paths []models.SitePath, target string) bool {
	for _, path := range paths {
		if path.Path == target {
			return true
		}
	}

	return false
}
