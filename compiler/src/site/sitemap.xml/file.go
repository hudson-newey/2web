package sitemapxml

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/xml"
	"hudson-newey/2web/src/models"
)

// TODO: refactor this so that the file writing is performed inside of the
// builder. We should aim for all writing to be performed inside of the builder.
func GenerateSitemap(paths []models.SitePath) {
	// The final newline is intentional so that the file is POSIX compliant.
	content := `<?xml version="1.0" encoding="UTF-8"?>
	<sitemap>
		<urlset>
			` + generateURLTags(paths) + `
		</urlset>
	</sitemap>
`

	file := xml.XmlFile{}
	file.AddContent(content)

	outputDirectory := cli.GetArgs().OutputPath
	if outputDirectory[len(outputDirectory)-1] != '/' {
		outputDirectory += "/"
	}
	outputFilePath := outputDirectory + "sitemap.xml"

	// In the 2web content model, pages are the entry points for web browsers.
	// Therefore, we re-use the same page writing logic to write out the
	// sitemap.xml file.
	pageModel := page.Page{
		Html: &file,
	}
	pageModel.WriteHtml(outputFilePath)
}

func generateURLTags(paths []models.SitePath) string {
	tags := ""
	for _, path := range paths {
		tags += "<url><loc>" + path.Path + "</loc></url>\n"
	}

	return tags
}
