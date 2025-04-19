package compiler

import (
	"fmt"
	"hudson-newey/2web/src/compiler/reactiveCompiler"
	"hudson-newey/2web/src/document/documentErrors"
	"hudson-newey/2web/src/lexer"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/models/reactiveEvent"
	"hudson-newey/2web/src/models/reactiveProperty"
	"hudson-newey/2web/src/models/reactiveVariable"
)

func Compile(filePath string, content string) string {
	reactiveVariables := []*models.ReactiveVariable{}

	// find all <script compiled></script> content
	compilerNodes := lexer.FindNodes[lexer.CompNode](content, compilerStartToken, compilerEndToken)

	propertyNodes := lexer.FindPropNodes[lexer.PropNode](content, propertyPrefix)
	eventNodes := lexer.FindPropNodes[lexer.EventNode](content, eventPrefix)

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

			reactiveVariables = append(reactiveVariables, &variableModel)
		}
	}

	for _, node := range propertyNodes {
		property, err := reactiveProperty.FromNode(node)
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
			if variable.Name == property.VarName {
				variable.AddProp(&property)
				property.BindVariable(variable)
				foundAssociatedProperty = true
				break
			}
		}

		if !foundAssociatedProperty {
			errorMessage := fmt.Sprintf("could not find compiler variable '%s' for property %s", property.VarName, property.Node.Selector)
			documentErrors.AddError(models.Error{
				FilePath: filePath,
				Message:  errorMessage,
			})
			continue
		}
	}

	for _, node := range eventNodes {
		event, err := reactiveEvent.FromNode(node)
		if err != nil {
			documentErrors.AddError(models.Error{
				FilePath: filePath,
				Message:  err.Error(),
			})
			continue
		}

		foundAssociatedProperty := false
		for _, variable := range reactiveVariables {
			if variable.Name == event.VarName {
				variable.AddEvent(&event)
				event.BindVariable(variable)
				foundAssociatedProperty = true
				break
			}
		}

		if !foundAssociatedProperty {
			errorMessage := fmt.Sprintf("could not find compiler variable '%s' for event %s", event.VarName, event.Node.Selector)
			documentErrors.AddError(models.Error{
				FilePath: filePath,
				Message:  errorMessage,
			})
			continue
		}
	}

	content = reactiveCompiler.CompileReactivity(filePath, content, reactiveVariables)

	return content
}
