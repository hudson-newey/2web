package javascript

import (
	"fmt"
	"hudson-newey/2web/src/constants"
	"sync"
)

const JsFunctionNamespace string = constants.CompilerNamespace + "func_"
const JsVarNamespace string = constants.CompilerNamespace + "var_"
const JsElementNamespace string = "data-" + constants.CompilerNamespace + "element_"

const ValueVar string = constants.CompilerNamespace + "value"

// nextNodeId is shared between all concurrently compiled pages because element
// names must be unique across the whole site.
//
// Access to it must be synchronized, otherwise two pages can be assigned the
// same id, producing colliding DOM selectors and JavaScript identifiers.
var nextNodeIdMutex sync.Mutex
var nextNodeId int = 0

// reserveNextNodeId returns the next unused node id.
func reserveNextNodeId() int {
	nextNodeIdMutex.Lock()
	defer nextNodeIdMutex.Unlock()

	id := nextNodeId
	nextNodeId++
	return id
}

func CreateJsFunctionName() string {
	functionName := fmt.Sprint(JsFunctionNamespace, reserveNextNodeId())
	return functionName
}

func CreateJsVariableName() string {
	variableName := fmt.Sprint(JsVarNamespace, reserveNextNodeId())
	return variableName
}

func CreateJsElementName() string {
	functionName := fmt.Sprint(JsElementNamespace, reserveNextNodeId())
	return functionName
}
