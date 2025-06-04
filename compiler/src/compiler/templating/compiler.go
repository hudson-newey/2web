package templating

import (
	"fmt"
	"hudson-newey/2web/src/compiler/lexer"
	"hudson-newey/2web/src/compiler/templating/imports"
	"hudson-newey/2web/src/compiler/templating/reactiveCompiler"
	"hudson-newey/2web/src/compiler/templating/textNode"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/page"
	"hudson-newey/2web/src/models"
	"hudson-newey/2web/src/models/reactiveEvent"
	"hudson-newey/2web/src/models/reactiveProperty"
	"hudson-newey/2web/src/models/reactiveVariable"
	"strings"
)

func Compile(filePath string, pageModel page.Page) page.Page {
	// "Mustache like" expressions e.g. {{ $count }} are a shorthand for an
	// element with only innerText.
	// Therefore, we expand all of the mustache expressions before finding
	// reactive property tokens
	pageModel.Html.Content = textNode.ExpandTextNodes(pageModel.Html.Content)

	importNodes := []lexer.LexNode[lexer.ImportNode]{}
	reactiveVariables := []*models.ReactiveVariable{}

	propertyNodes := lexer.FindPropNodes[lexer.PropNode](pageModel.Html.Content, propertyPrefix)
	eventNodes := lexer.FindPropNodes[lexer.EventNode](pageModel.Html.Content, eventPrefix)

	for _, node := range pageModel.JavaScript {
		// At the moment, we do not support mixing compiler and runtime scripts
		if !node.IsCompilerOnly() {
			continue
		}

		commentNodes := lexer.FindNodes[lexer.LineCommentNode](node.Content, lineCommentStartToken, newLineToken)
		commentLessContent := node.Content
		for _, commentNode := range commentNodes {
			commentLessContent = strings.ReplaceAll(commentLessContent, commentNode.Selector, "")
		}

		variableNodes := lexer.FindNodes[lexer.VarNode](commentLessContent, variableToken, statementEndToken)

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

		newImportNodes := lexer.FindNodes[lexer.ImportNode](node.Content, importPrefix, statementEndToken)
		importNodes = append(importNodes, newImportNodes...)
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

	pageModel.Html.Content = imports.EvaluateImports(filePath, pageModel.Html.Content, importNodes)

	pageModel = reactiveCompiler.CompileReactivity(filePath, pageModel, reactiveVariables)

	return pageModel
}
