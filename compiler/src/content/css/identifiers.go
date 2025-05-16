package css

import (
	"fmt"
	"hudson-newey/2web/src/constants"
)

var nextNodeId int = 0

func CreateCssId() string {
	idName := fmt.Sprint(constants.CompilerNamespace+"id_", nextNodeId)
	nextNodeId++
	return idName
}
