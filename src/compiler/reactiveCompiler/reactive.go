package reactiveCompiler

import (
	"fmt"
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/javascript"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileReactiveVar(
	content string,
	varNode *models.ReactiveVariable,
) string {
	callbackName := javascript.CreateJsFunctionName()
	variableName := javascript.CreateJsVariableName()

	htmlSource := fmt.Sprintf(`
    <script>
      let %s = %s;
      function %s(%s) {
        {{range .Props}}
            document.querySelector("[{{.Node.Selector}}]")[{{.PropName}}] = %s;
        {{end}}
      }
    </script>
  `, variableName, varNode.InitialValue, callbackName, javascript.ValueVar, javascript.ValueVar)

	injectableTemplate, err := document.BuildTemplate(htmlSource, *varNode)
	if err != nil {
		panic(err)
	}

	// TODO: remove this stupid hack once the reactive compiler doesn't rely
	// on the static compiler to make element references
	injectableTemplate = strings.ReplaceAll(injectableTemplate, "\\u0022", "\"")

	content = document.InjectContent(content, injectableTemplate, document.Body)

	for _, event := range varNode.Events {
		reactiveReducer := strings.ReplaceAll(event.Reducer, varNode.Name, variableName)

		eventBindingAttribute := fmt.Sprintf("on%s='%s(%s)'", event.EventName, callbackName, reactiveReducer)
		content = strings.ReplaceAll(content, event.Node.Selector, eventBindingAttribute)
	}

	return content + htmlSource
}
