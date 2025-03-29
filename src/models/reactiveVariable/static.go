package reactiveVariable

import (
	"fmt"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
)

func FromNode(node lexer.LexNode[lexer.VarNode]) (models.ReactiveVariable, error) {
	if len(node.Tokens) != 3 || node.Tokens[1] != "=" {
		return models.ReactiveVariable{}, fmt.Errorf("Incorrect compiler variable format:\nUsage: $ variableName = 'variableValue'")
	}

	varName := node.Tokens[0]
	varValue := node.Tokens[2]

	variableModel := models.ReactiveVariable{
		Name:         "$" + varName,
		InitialValue: varValue,
	}

	return variableModel, nil
}
