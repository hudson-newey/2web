package compiler

import (
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
)

func Compile(path string, content string) string {
	compilerNodes := lexer.FindNodes(content, compilerStartToken, compilerEndToken)

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

			println(varName, varValue)
		}
	}

	return content
}
