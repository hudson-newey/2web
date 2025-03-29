package compiler

import (
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/models/reactiveProperty"
	"hudson-newey/2web/src/models/reactiveVariable"
)

func Compile(path string, content string) string {
	compilerNodes := lexer.FindNodes[lexer.CompNode](content, compilerStartToken, compilerEndToken)

	reactiveVariables := []models.ReactiveVariable{}
	reactiveProperties := []models.ReactiveProperty{}

	propertyNodes := lexer.FindNodes[lexer.PropNode](content, reactiveStartToken, reactiveEndToken)

	for _, propertyNode := range propertyNodes {
		propertyModel, err := reactiveProperty.FromNode(propertyNode)
		if err != nil {
			documentErrors.AddError(models.Error{
				FilePath: path,
				Message:  err.Error(),
			})
			continue
		}

		reactiveProperties = append(reactiveProperties, propertyModel)
	}

	for _, node := range compilerNodes {
		variableNodes := lexer.FindNodes[lexer.VarNode](node.Content, variableToken, statementEndToken)

		for _, variableNode := range variableNodes {
			variableModel, err := reactiveVariable.FromNode(variableNode)
			if err != nil {
				documentErrors.AddError(models.Error{
					FilePath: path,
					Message:  err.Error(),
				})
				continue
			}

			reactiveVariables = append(reactiveVariables, variableModel)
		}

	}

	content = compileReactivity(content, reactiveVariables, reactiveProperties)

	return content
}
