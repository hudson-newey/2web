package reactiveCompiler

import (
	"fmt"
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/javascript"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileAssignmentProp(content string, varNode *models.ReactiveVariable) string {
	callbackName := javascript.CreateJsFunctionName()

	htmlSource := `
    <script>
      function ` + callbackName + `(newValue) {
        {{range .Props}}
            document.querySelector("[{{.Node.Selector}}]")[{{.PropName}}] = newValue;
        {{end}}
      }
    </script>
	`

	injectableTemplate, err := document.BuildTemplate(htmlSource, *varNode)
	if err != nil {
		panic(err)
	}

	content = document.InjectContent(content, injectableTemplate, document.Body)

	for _, event := range varNode.Events {
		eventBindingAttribute := fmt.Sprintf("on%s='%s(\"%s\")'", event.EventName, callbackName, event.Reducer)
		content = strings.ReplaceAll(content, event.Node.Selector, eventBindingAttribute)
	}

	return content
}
