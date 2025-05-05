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
	elementSelector := javascript.CreateJsElementName()

	functionContent := ""
	for _, propNode := range varNode.Props {
		content = strings.ReplaceAll(content, propNode.Node.Selector, propNode.Node.Selector+" "+elementSelector)
		functionContent += fmt.Sprintf(`
      document.querySelector("[%s]")["%s"] = %s;
    `, elementSelector, propNode.PropName, javascript.ValueVar)
	}

	updateJsSource := fmt.Sprintf(`
    <script>
      let %s = %s;
      function %s(%s) {
        %s
      }
    </script>
  `, variableName, varNode.InitialValue, callbackName, javascript.ValueVar, functionContent)

	injectableTemplate, err := document.BuildTemplate(updateJsSource, *varNode)
	if err != nil {
		panic(err)
	}

	// TODO: remove this stupid hack once the reactive compiler doesn't rely
	// on the static compiler to make element references
	injectableTemplate = strings.ReplaceAll(injectableTemplate, "\\u0022", "\"")

	content = document.InjectContent(content, injectableTemplate, document.Body)

	for _, event := range varNode.Events {
		reactiveReducer := strings.ReplaceAll(event.Reducer, varNode.Name, variableName)

		// e.g. <button onclick="count = count + 1; updateCount(count)">Increment</button>
		eventBindingAttribute := ""
		if useDoubleQuotes(event.Reducer) {
			eventBindingAttribute =
				fmt.Sprintf(
					"on%s=\"%s = %s; %s(%s)\"",
					event.EventName,
					variableName,
					reactiveReducer,
					callbackName,
					variableName,
				)
		} else {
			eventBindingAttribute =
				fmt.Sprintf(
					"on%s='%s = %s; %s(%s)'",
					event.EventName,
					variableName,
					reactiveReducer,
					callbackName,
					variableName,
				)
		}

		content = strings.ReplaceAll(content, event.Node.Selector, eventBindingAttribute)
	}

	return content
}
