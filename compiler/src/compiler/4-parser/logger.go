package parser

import (
	"fmt"
	"hudson-newey/2web/src/compiler/4-parser/ast"
)

func PrintVerboseParser(structure ast.AbstractSyntaxTree) {
	recursiveVerbosePrint(structure, 0)
}

func recursiveVerbosePrint(structure ast.AbstractSyntaxTree, depth int) {
	for _, node := range structure {
		fmt.Println(node.Type())

		recursiveVerbosePrint(node.Children(), depth+1)
	}
}
