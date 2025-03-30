package compiler

import (
	"fmt"
	"hudson-newey/2web/src/document"
	"hudson-newey/2web/src/models"
	"strings"
)

// TODO: move this to another place
var nextNodeId int = 0

func compileReactivity(
	filePath string,
	content string,
	varNodes []*models.ReactiveVariable,
) string {
	for _, varNode := range varNodes {
		if varNode.Type() >= models.Assignment {
			content = compileAssignmentProp(content, varNode)
		}

		if varNode.Type() >= models.StaticProperty {
			content = compileStaticProperty(content, varNode)
		}

		content = compileStatic(content, varNode)
	}

	return content
}

func compileStatic(content string, varNode *models.ReactiveVariable) string {
	content = strings.ReplaceAll(content, mustacheStartToken+varNode.Name+mustacheEndToken, varNode.InitialValue)
	content = strings.ReplaceAll(content, varNode.Name, varNode.InitialValue)

	return content
}

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

func compileAssignmentProp(content string, varNode *models.ReactiveVariable) string {
	// callbackName := fmt.Sprint("__2_", nextNodeId)
	// elementSelector := fmt.Sprint("data-2='", nextNodeId, "'")

	// htmlSource := `
	//     <script>
	//         function ` + callbackName + `(newValue) {
	//             {{range .Props}}
	//                 document.querySelector("[data-2='{{.Node.Selector}}']").{{.PropName}} = newValue;
	//             {{end}}
	//         }
	//     </script>
	// `

	for _, event := range varNode.Events {
		content = strings.ReplaceAll(content, event.Node.Selector, "")
	}

	return content
}
