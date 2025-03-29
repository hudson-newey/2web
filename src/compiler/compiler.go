package compiler

import (
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
	"strings"
)

func Compile(path string, content string) string {
	compilerNodes := lexer.FindNodes(content, compilerStartToken, compilerEndToken)

	reactiveVariables := []models.ReactiveVariable{}

	for _, node := range compilerNodes {
		variableNodes := lexer.FindNodes(node.Content, variableToken, statementEndToken)
		for _, variableNode := range variableNodes {
			if len(variableNode.Tokens) != 3 || variableNode.Tokens[1] != variableAssignmentToken {
				documentErrors.AddError(models.Error{
					FilePath: path,
					Message:  "Incorrect compiler variable format:\nUsage: $ variableName = 'variableValue'",
				})
				continue
			}

			varName := variableNode.Tokens[0]
			varValue := variableNode.Tokens[2]

			variableModel := models.ReactiveVariable{
				Name:         "$" + varName,
				InitialValue: varValue,
			}

			reactiveVariables = append(reactiveVariables, variableModel)
		}
	}

	for _, variable := range reactiveVariables {
		content = strings.ReplaceAll(content, mustacheStartToken+variable.Name+mustacheEndToken, variable.InitialValue)
		content = strings.ReplaceAll(content, variable.Name, variable.InitialValue)
	}

	return content
}
