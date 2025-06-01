package templates

import (
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

const serverEntryContent string = `import { runServer } from "@two-web/kit/ssr";

runServer();
`

func SsrTemplate() {
	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("src/server/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        "server/ssr.ts",
			Content:     serverEntryContent,
			IsDirectory: false,
		},
	}

	files.WriteFiles(templateFiles)
}
