package reactiveCompiler

import (
	"fmt"
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileStaticProperty(content string, varNode *models.ReactiveVariable) string {
	for _, propNode := range varNode.Props {
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
