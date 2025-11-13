package page

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/filesystem"
)

func (model *Page) WriteHtml(filePath string) {
	write(model.Html.Content, filePath)
	model.WriteAssets()
}

// Writes assets like CSS and JavaScript to their lazy loaded modules
func (model *Page) WriteAssets() {
	for _, file := range model.Css {
		write(file.RawContent(), file.OutputPath())
	}

	for _, file := range model.JavaScript {
		write(file.RawContent(), file.OutputPath())
	}

	for _, file := range model.Assets {
		derefFile := *file
		writeBinary(derefFile.Data(), derefFile.OutputPath())
	}
}

func write(content string, outputPath string) {
	// Make sure that the files content ends in a newline so that it is POSIX
	// compliant.
	// This is needed otherwise the non-blocking write may not write the correct
	// content because it uses low-level syscalls.
	content += "\n"

	if *cli.GetArgs().ToStdout {
		fmt.Println(content)
	} else {
		filesystem.WriteFile([]byte(content), outputPath)
	}
}

func writeBinary(content []byte, outputPath string) {
	if *cli.GetArgs().ToStdout {
		fmt.Println(string(content))
	} else {
		filesystem.WriteFile(content, outputPath)
	}
}
