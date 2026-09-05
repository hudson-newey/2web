package robotstxt

import (
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/content/txt"
)

// The final newline is intentional so that the file is POSIX compliant.
const defaultRobotsTxtContent = `User-agent: *
Disallow:
`

func GenerateRobotsTxt() {
	file := txt.TxtFile{}
	file.AddContent(defaultRobotsTxtContent)

	outputDirectory := cli.GetArgs().OutputPath
	if outputDirectory[len(outputDirectory)-1] != '/' {
		outputDirectory += "/"
	}
	outputFilePath := outputDirectory + "robots.txt"

	// In the 2web content model, pages are the entry points for web browsers.
	// Therefore, we re-use the same page writing logic to write out the
	// robots.txt file.
	pageModel := page.Page{
		Html: &file,
	}
	pageModel.WriteHtml(outputFilePath)
}
