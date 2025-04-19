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

	htmlSource := fmt.Sprintf(`
    <script>
      function %s(newValue) {
        {{range .Props}}
            document.querySelector("[{{.Node.Selector}}]")[{{.PropName}}] = newValue;
        {{end}}
      }
    </script>
	`, callbackName)

	injectableTemplate, err := document.BuildTemplate(htmlSource, *varNode)
	if err != nil {
		panic(err)
	}

	// TODO: remove this stupid hack once the reactive compiler doesn't rely
	// on the static compiler to make element references
	injectableTemplate = strings.ReplaceAll(injectableTemplate, "\\u0022", "\"")

	content = document.InjectContent(content, injectableTemplate, document.Body)

	for _, event := range varNode.Events {
		eventBindingAttribute := ""
		if useDoubleQuotes(event.Reducer) {
			eventBindingAttribute = fmt.Sprintf("on%s=\"%s(%s)\"", event.EventName, callbackName, event.Reducer)
		} else {
			eventBindingAttribute = fmt.Sprintf("on%s='%s(%s)'", event.EventName, callbackName, event.Reducer)
		}

		content = strings.ReplaceAll(content, event.Node.Selector, eventBindingAttribute)
	}

	return content
}
