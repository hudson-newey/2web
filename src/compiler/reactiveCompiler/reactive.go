package reactiveCompiler

import (
	"hudson-newey/2web/src/javascript"
	"hudson-newey/2web/src/models"
)

func compileReactiveTemplate(
	content string,
	varNode *models.ReactiveVariable,
) string {
	callbackName := javascript.CreateJsFunctionName()
	variableName := javascript.CreateJsVariableName()

	htmlSource := `
    <script>
      let ` + variableName + ` = 0;
      function ` + callbackName + `(newValue) {
        {{range .Props}}
            document.querySelector("[{{.Node.Selector}}]")[{{.PropName}}] = newValue;
        {{end}}
      }
    </script>
  `

	return content + htmlSource
}
