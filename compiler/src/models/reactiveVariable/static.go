package reactiveVariable

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/models"
)

func FromNode(node lexer.LexNode[lexer.VarNode]) (models.ReactiveVariable, error) {
	if len(node.Tokens) != 3 || node.Tokens[1] != "=" {
		errorMessage := fmt.Errorf("incorrect reactive variable assignment:\n\tExpected: $ variableName = variableValue\n\tFound: %s", node.Selector)
		return models.ReactiveVariable{}, errorMessage
	}

	varName := node.Tokens[0]
	varValue := node.Tokens[2]

	variableModel := models.ReactiveVariable{
		Node:         &node,
		Name:         "$" + varName,
		InitialValue: varValue,
		Reactive:     false,
	}

	return variableModel, nil
}
