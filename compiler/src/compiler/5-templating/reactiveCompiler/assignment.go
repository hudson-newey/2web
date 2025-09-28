package reactiveCompiler

import (
	"fmt"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileAssignmentVar(content string, varNode *models.ReactiveVariable) string {
	callbackName := javascript.CreateJsFunctionName()

	uniquePropSelectors := getUniqueSelectors(varNode.Props)

	functionContent := ""
	for _, propNode := range uniquePropSelectors {
		elementSelector := javascript.CreateJsElementName()
		content = strings.ReplaceAll(content, propNode.Node.Selector, propNode.Node.Selector+" "+elementSelector)

		selectorCount := strings.Count(content, propNode.Node.Selector)
		if selectorCount > 1 {
			functionContent += fmt.Sprintf(
				`document.querySelectorAll("[%s]").forEach((__2_element_ref_mod) => __2_element_ref_mod["%s"] = %s);`,
				elementSelector, propNode.PropName, javascript.ValueVar,
			)
		} else {
			functionContent += fmt.Sprintf(
				`document.querySelector("[%s]")["%s"] = %s;`,
				elementSelector, propNode.PropName, javascript.ValueVar,
			)
		}
	}

	// There's a newline at the start of this script tag so that when it is
	// appended to the body, it's on its own line, and semantically distinct.
	updateJsSource := fmt.Sprintf(`
		<script type="module">
      globalThis.%s = (%s) => {
        %s
      }
    </script>`,
		callbackName, javascript.ValueVar, functionContent,
	)

	injectableTemplate, err := document.BuildTemplate(updateJsSource, *varNode)
	if err != nil {
		panic(err)
	}

	// TODO: remove this stupid hack once the reactive compiler doesn't rely
	// on the static compiler to make element references
	injectableTemplate = strings.ReplaceAll(injectableTemplate, "\\u0022", "\"")

	content = document.InjectContent(content, injectableTemplate, document.BodyTop)

	for _, event := range varNode.Events {
		eventBindingAttribute := ""
		if UseDoubleQuotes(event.Reducer) {
			eventBindingAttribute = fmt.Sprintf("on%s=\"%s(%s)\"", event.EventName, callbackName, event.Reducer)
		} else {
			eventBindingAttribute = fmt.Sprintf("on%s='%s(%s)'", event.EventName, callbackName, event.Reducer)
		}

		content = strings.ReplaceAll(content, event.Node.Selector, eventBindingAttribute)
	}

	return content
}
