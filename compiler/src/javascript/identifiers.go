package javascript

import "fmt"

const ValueVar string = "__2_value"

var nextNodeId int = 0

func CreateJsFunctionName() string {
	functionName := fmt.Sprint("__2_func_", nextNodeId)
	nextNodeId++
	return functionName
}

func CreateJsVariableName() string {
	variableName := fmt.Sprint("__2_var_", nextNodeId)
	nextNodeId++
	return variableName
}

func CreateJsElementName() string {
	functionName := fmt.Sprint("data-__2_element_", nextNodeId)
	nextNodeId++
	return functionName
}
