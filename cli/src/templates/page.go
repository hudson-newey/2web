package templates

import (
	"fmt"
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

func PageTemplate(name string) {
	dirPath := "src/"

	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        dirPath + name + ".html",
			IsDirectory: false,
			Content:     createPageContent(name),
		},
	}

	files.WriteFiles(templateFiles)
}

func createPageContent(name string) string {
	return fmt.Sprintf(`<h1>%s works!</h1>`, name)
}
