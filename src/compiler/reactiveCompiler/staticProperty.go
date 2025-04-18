package reactiveCompiler

import (
	"fmt"
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/javascript"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileStaticProperty(content string, varNode *models.ReactiveVariable) string {
	for _, propNode := range varNode.Props {
		elementSelector := javascript.CreateJsElement()
		content = strings.ReplaceAll(content, propNode.Node.Selector, elementSelector)

		// we use the square brackets here because some properties have dashes which
		// cannot be acceded with a period
		htmlSource := fmt.Sprintf(`
      <script>
          document.querySelector("[%s]")["%s"] = {{.Variable.InitialValue}};
      </script>
    `, elementSelector, propNode.PropName)

		injectableTemplate, err := document.BuildTemplate(htmlSource, *propNode)
		if err != nil {
			panic(err)
		}

		content = document.InjectContent(content, injectableTemplate, document.Body)
	}

	return content
}
