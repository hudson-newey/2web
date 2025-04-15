package javascript

import "fmt"

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
