package page

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	"hudson-newey/2web/src/utils"
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
		if file.IsCompilerOnly() {
			continue
		}

		write(file.RawContent(), file.OutputPath())
	}

	for _, file := range model.Assets {
		derefFile := *file
		writeBinary(derefFile.Data(), derefFile.OutputPath())
	}
}

func write(content string, outputPath string) {
	if *cli.GetArgs().ToStdout {
		fmt.Println(content)
	} else {
		utils.WriteFile(content, outputPath)
	}
}

func writeBinary(content []byte, outputPath string) {
	if *cli.GetArgs().ToStdout {
		fmt.Println(string(content))
	} else {
		utils.WriteBinaryFile(content, outputPath)
	}
}
