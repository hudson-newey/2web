package imports

import (
	"fmt"
	"hudson-newey/2web/src/compiler/pageBuilder"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"os"
	"strings"
)

func ExpandComponentImports(
	workingPath string,
	content string,
	components []*models.Component,
) string {
	result := content

	for _, component := range components {
		result = expandImport(workingPath, content, component)
	}

	return result
}

func expandImport(
	workingPath string,
	content string,
	component *models.Component,
) string {
	pageModel, err := buildComponent(component)
	if err != nil {
		documentErrors.AddError(models.Error{
			FilePath: workingPath,
			Message:  err.Error(),
		})

		// We return the input content without modification in the hopes that the
		// page will be semi-functional and assist in debugging how the error
		// occurred.
		return content
	}

	templateSelector := fmt.Sprintf("<%s />", component.DomSelector)
	result := strings.ReplaceAll(content, templateSelector, pageModel.Html.Content)

	// The HTML model might contain links to lazy loaded assets like CSS and
	// JavaScript.
	// By writing the assets, these separate lazy loaded files will be written.
	// For simplicity sake, you can think of this writing the JavaScript and CSS
	// for the component, but not writing an associated .html file.
	pageModel.WriteAssets()

	return result
}

func buildComponent(component *models.Component) (page.Page, error) {
	inputPath, err := component.ComponentPath()
	if err != nil {
		return page.Page{}, err
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return page.Page{}, err
	}

	// We don't want to expand HTML partials for component imports
	fullDocumentContent := ""
	if markdown.IsMarkdownFile(inputPath) {
		markdownFile := markdown.MarkdownFile{
			Content: string(data),
		}
		fullDocumentContent = markdownFile.ToHtml().Content
	} else {
		fullDocumentContent = string(data)
	}

	pageModel := pageBuilder.BuildPage(fullDocumentContent)

	return pageModel, nil
}
