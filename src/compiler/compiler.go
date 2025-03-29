package compiler

import (
	"fmt"
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/models/reactiveProperty"
	"hudson-newey/2web/src/models/reactiveVariable"
)

func Compile(filePath string, content string) string {
	compilerNodes := lexer.FindNodes[lexer.CompNode](content, compilerStartToken, compilerEndToken)

	reactiveVariables := []models.ReactiveVariable{}
	reactiveProperties := []models.ReactiveProperty{}

	propertyNodes := lexer.FindNodes[lexer.PropNode](content, reactiveStartToken, reactiveEndToken)

	for _, node := range compilerNodes {
		variableNodes := lexer.FindNodes[lexer.VarNode](node.Content, variableToken, statementEndToken)

		for _, variableNode := range variableNodes {
			variableModel, err := reactiveVariable.FromNode(variableNode)
			if err != nil {
				documentErrors.AddError(models.Error{
					FilePath: filePath,
					Message:  err.Error(),
				})
				continue
			}

			reactiveVariables = append(reactiveVariables, variableModel)
		}
	}

	for _, propertyNode := range propertyNodes {
		propertyModel, err := reactiveProperty.FromNode(propertyNode)
		if err != nil {
			documentErrors.AddError(models.Error{
				FilePath: filePath,
				Message:  err.Error(),
			})
			continue
		}

		// find the reactive variable that matches the binding name
		foundAssociatedProperty := false
		for _, variable := range reactiveVariables {
			if variable.Name == propertyModel.BindingName {
				variable.Type = models.Assignment

				propertyModel.AddBinding(&variable)

				foundAssociatedProperty = true
				break
			}
		}

		if !foundAssociatedProperty {
			errorMessage := fmt.Sprintf("could not find reactive variable '%s' for property %s", propertyModel.BindingName, propertyModel.Node.Selector)
			documentErrors.AddError(models.Error{
				FilePath: filePath,
				Message:  errorMessage,
			})
			continue
		}

		reactiveProperties = append(reactiveProperties, propertyModel)
	}

	content = compileReactivity(filePath, content, reactiveVariables, reactiveProperties)

	return content
}
