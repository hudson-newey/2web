package reactiveCompiler

import (
	"fmt"
	"hudson-newey/2web/src/content/document"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"strings"
)

func compileReactiveVar(
	pageModel *page.Page,
	varNode *models.ReactiveVariable,
) {
	callbackName := javascript.CreateJsFunctionName()
	variableName := javascript.CreateJsVariableName()

	uniquePropSelectors := getUniqueSelectors(varNode.Props)

	content := pageModel.Html.Content

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
      globalThis.%s = %s;
      globalThis.%s = (%s) => {
        %s
      }
    </script>`,
		variableName, varNode.InitialValue, callbackName, javascript.ValueVar, functionContent,
	)

	injectableTemplate, err := document.BuildTemplate(updateJsSource, *varNode)
	if err != nil {
		panic(err)
	}

	// TODO: remove this stupid hack once the reactive compiler doesn't rely
	// on the static compiler to make element references
	injectableTemplate = strings.ReplaceAll(injectableTemplate, "\\u0022", "\"")

	// If there is no body tag, we just append to the end of the content.
	// TODO: This should be improved once we start using lazy loaded scripts
	// instead of inlining reactivity as scripts in HTML content.
	if document.HasBodyTag(content) {
		content = document.InjectContent(content, injectableTemplate, document.BodyTop)
	} else {
		content = document.InjectContent(content, injectableTemplate, document.Leading)
	}

	for _, event := range varNode.Events {
		reactiveReducer := strings.ReplaceAll(event.Reducer, varNode.Name, variableName)

		// e.g. <button onclick="count = count + 1; updateCount(count)">Increment</button>
		eventBindingAttribute := ""
		if UseDoubleQuotes(event.Reducer) {
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

	pageModel.Html.Content = content
}
