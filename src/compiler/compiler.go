package compiler

import (
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/models/reactiveVariable"
)

func Compile(path string, content string) string {
	compilerNodes := lexer.FindNodes[lexer.CompNode](content, compilerStartToken, compilerEndToken)

	reactiveVariables := []models.ReactiveVariable{}
	reactiveProperties := []models.ReactiveProperty{}

	for _, node := range compilerNodes {
		variableNodes := lexer.FindNodes[lexer.VarNode](node.Content, variableToken, statementEndToken)
		propertyNodes := lexer.FindNodes[lexer.PropNode](node.Content, reactiveStartToken, reactiveEndToken)

		for _, variableNode := range variableNodes {
			variableModel, err := reactiveVariable.FromNode(variableNode)
			if err != nil {
				documentErrors.AddError(models.Error{
					FilePath: path,
					Message:  "Incorrect compiler variable format:\nUsage: $ variableName = 'variableValue'",
				})
				continue
			}

			reactiveVariables = append(reactiveVariables, variableModel)
		}
	}

	content = compileReactivity(content, reactiveVariables, reactiveProperties)

	return content
}
