package templating

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/markdown"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"os"
	"strings"
)

func ExpandComponentImports(
	workingPath string,
	content page.Page,
	components []*models.Component,
) page.Page {
	result := content

	for _, component := range components {
		result = expandImport(workingPath, content, component)
	}

	return result
}

func expandImport(
	workingPath string,
	page page.Page,
	component *models.Component,
) page.Page {
	componentModel, err := buildComponent(component)
	if err != nil {
		documentErrors.AddErrors(
			models.NewError(
				err.Error(),
				workingPath,
				lexer.Position{},
			),
		)

		// We return the input content without modification in the hopes that the
		// page will be semi-functional and assist in debugging how the error
		// occurred.
		return page
	}

	templateSelector := fmt.Sprintf("<%s />", component.DomSelector)
	page.Html.Content = strings.ReplaceAll(
		page.Html.Content,
		templateSelector,
		componentModel.Html.Content,
	)

	// Notice that for each component, its styles and scripts are only added to
	// the page once.
	for _, cssFile := range componentModel.Css {
		page.AddStyle(cssFile)
	}

	for _, jsFile := range componentModel.JavaScript {
		page.AddScript(jsFile)
	}

	// The HTML model might contain links to lazy loaded assets like CSS and
	// JavaScript.
	// By writing the assets, these separate lazy loaded files will be written.
	// For simplicity sake, you can think of this writing the JavaScript and CSS
	// for the component, but not writing an associated .html file.
	componentModel.WriteAssets()

	return page
}

func buildComponent(component *models.Component) (page.Page, error) {
	inputPath, err := component.ComponentPath()
	if err != nil {
		return page.NewPage(), err
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return page.NewPage(), err
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

	pageModel := BuildPage(inputPath, fullDocumentContent)

	return pageModel, nil
}
