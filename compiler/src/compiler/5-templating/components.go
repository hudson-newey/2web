package templating

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"strings"
)

func expandComponentImports(
	workingPath string,
	pageModel *page.Page,
	components []*models.Component,
) {
	for _, component := range components {
		expandImport(workingPath, pageModel, component)
	}
}

func expandImport(
	workingPath string,
	pageModel *page.Page,
	component *models.Component,
) *page.Page {
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
		return pageModel
	}

	templateSelector := fmt.Sprintf("<%s />", component.DomSelector)
	pageModel.Html.Content = strings.ReplaceAll(
		pageModel.Html.Content,
		templateSelector,
		componentModel.Html.Content,
	)

	// Notice that for each component, its styles and scripts are only added to
	// the page once.
	for _, cssFile := range componentModel.Css {
		pageModel.AddStyle(cssFile)
	}

	for _, jsFile := range componentModel.JavaScript {
		pageModel.AddScript(jsFile)
	}

	for _, twoScript := range componentModel.TwoScript {
		pageModel.AddTwoScript(twoScript)
	}

	for _, asset := range componentModel.Assets {
		pageModel.AddAsset(asset)
	}

	// The HTML model might contain links to lazy loaded assets like CSS and
	// JavaScript.
	// By writing the assets, these separate lazy loaded files will be written.
	// For simplicity sake, you can think of this writing the JavaScript and CSS
	// for the component, but not writing an associated .html file.
	componentModel.WriteAssets()

	return pageModel
}

func buildComponent(component *models.Component) (page.Page, error) {
	inputPath, err := component.ComponentPath()
	if err != nil {
		return page.NewPage(), err
	}

	// Use the injected build function to avoid an import cycle between
	// the templating and builder packages. This is set by the builder package
	// at init-time. If it's not set, this will return (empty, false).
	pageModel, success := BuildComponentPage(inputPath, false)
	if !success {
		return page.NewPage(), fmt.Errorf(
			"failed to build component at path: %s",
			inputPath,
		)
	}

	return pageModel, nil
}
