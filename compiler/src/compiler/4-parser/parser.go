package parser

import lexer "hudson-newey/2web/src/compiler/2-lexer"

type AbstractSyntaxTree []Node

func CreateAst(structure []lexer.V2LexNode) AbstractSyntaxTree {
	var ast AbstractSyntaxTree
	for _, node := range structure {
		ast = append(ast, processNode(node))
	}

	return ast
}

func processNode(structure lexer.V2LexNode) Node {
	// TODO: During development, I just want to treat everything as a text node
	return textNode(structure)
}

func textNode(structure lexer.V2LexNode) Node {
	return &MarkupTextNode{
		content: structure.Content,
	}
}
