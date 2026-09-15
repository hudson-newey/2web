package templates

import (
	"os"

	"github.com/hudson-newey/2web-cli/src/files"
)

const serverEntryContent string = `import { runServer } from "@two-web/kit/ssr";

// Runs the ssr server: compiled html documents are server side rendered,
// compiled assets are served statically, and the compiled server routes
// (.server.ts files) are mounted.
//
// The server listens on http://localhost:5173. Change the port with
// "2web serve --port <port>" or the PORT environment variable.
runServer();
`

func SsrTemplate() {
	// ignore errors from this because we expect this to fail (because the
	// directory) already exists
	os.Mkdir("server/", os.ModePerm)

	templateFiles := []files.File{
		{
			Path:        "server/ssr.ts",
			Content:     serverEntryContent,
			IsDirectory: false,
		},
	}

	files.WriteFiles(templateFiles)
}
