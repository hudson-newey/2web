package devtools

import (
	"hudson-newey/2web/src/content/document"
)

func InjectDevTools(pageContent string) string {
	devtoolsTemplateResult := createDevtoolsTemplate()
	return document.InjectContent(pageContent, devtoolsTemplateResult, document.BodyTop)
}

func createDevtoolsTemplate() string {
	devtoolsHtml, err := document.BuildTemplate(devtoolsHtmlSource(), nil)
	if err != nil {
		panic(err)
	}

	return devtoolsHtml
}
