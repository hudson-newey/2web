package reactiveProperty

import (
	"fmt"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
)

func FromNode(node lexer.LexNode[lexer.PropNode]) (models.ReactiveProperty, error) {
	if len(node.Tokens) != 2 {
		errorMessage := fmt.Errorf("incorrect reactive property assignment:\n\tExpected: [propertyName $variable]\n\tFound: %s", node.Selector)
		return models.ReactiveProperty{}, errorMessage
	}

	propName := node.Tokens[0]
	bindingName := node.Tokens[1]

	propertyModel := models.ReactiveProperty{
		Node:     &node,
		PropName: propName,
		VarName:  bindingName,
	}

	return propertyModel, nil
}
