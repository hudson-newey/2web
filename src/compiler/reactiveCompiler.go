package compiler

import (
	"fmt"
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileReactivity(
	filePath string,
	content string,
	varNodes []*models.ReactiveVariable,
	propNodes []*models.ReactiveProperty,
) string {
	for _, varNode := range varNodes {
		switch varNode.Type() {
		case models.Static:
			content = compileStatic(content, varNode)
		case models.StaticProperty:
			content = compileStaticProperty(content, varNode)
		}
	}

	return content
}

func compileStatic(content string, varNode *models.ReactiveVariable) string {
	content = strings.ReplaceAll(content, mustacheStartToken+varNode.Name+mustacheEndToken, varNode.InitialValue)
	content = strings.ReplaceAll(content, varNode.Name, varNode.InitialValue)

	return content
}

func compileStaticProperty(content string, varNode *models.ReactiveVariable) string {
	nextNodeId := 0
	for _, propNode := range varNode.Bindings {
		elementSelector := fmt.Sprint("data-2='", nextNodeId, "'")
		content = strings.ReplaceAll(content, propNode.Node.Selector, elementSelector)
		nextNodeId++

		htmlSource := `
            <script>
                document.querySelector("[` + elementSelector + `]").` + propNode.PropName + ` = {{.Variable.InitialValue}};
            </script>
        `

		injectableTemplate, err := document.BuildTemplate(htmlSource, *propNode)
		if err != nil {
			panic(err)
		}

		content = document.InjectContent(content, injectableTemplate, document.Body)
	}

	return content
}
