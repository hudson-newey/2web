package javascript

import (
	"fmt"
	"hudson-newey/2web/src/constants"
)

const JsFunctionNamespace string = constants.CompilerNamespace + "func_"
const JsVarNamespace string = constants.CompilerNamespace + "var_"
const JsElementNamespace string = "data-" + constants.CompilerNamespace + "element_"

const ValueVar string = constants.CompilerNamespace + "value"

var nextNodeId int = 0

func CreateJsFunctionName() string {
	functionName := fmt.Sprint(JsFunctionNamespace, nextNodeId)
	nextNodeId++
	return functionName
}

func CreateJsVariableName() string {
	variableName := fmt.Sprint(JsVarNamespace, nextNodeId)
	nextNodeId++
	return variableName
}

func CreateJsElementName() string {
	functionName := fmt.Sprint(JsElementNamespace, nextNodeId)
	nextNodeId++
	return functionName
}
