package javascript

import (
	"fmt"
	"hudson-newey/2web/src/constants"
)

const ValueVar string = constants.CompilerNamespace + "value"

var nextNodeId int = 0

func CreateJsFunctionName() string {
	functionName := fmt.Sprint(constants.CompilerNamespace+"func_", nextNodeId)
	nextNodeId++
	return functionName
}

func CreateJsVariableName() string {
	variableName := fmt.Sprint(constants.CompilerNamespace+"var_", nextNodeId)
	nextNodeId++
	return variableName
}

func CreateJsElementName() string {
	functionName := fmt.Sprint("data-"+constants.CompilerNamespace+"element_", nextNodeId)
	nextNodeId++
	return functionName
}
