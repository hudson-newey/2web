package parser

import (
	"fmt"
	"hudson-newey/2web/src/compiler/4-parser/ast"
)

func PrintVerboseParser(structure ast.AbstractSyntaxTree) {
	recursiveVerbosePrint(structure, 0)

	// Draw an empty line at the end for better readability and to visually
	// separate this log section from any subsequent logs.
	fmt.Println()
}

func recursiveVerbosePrint(structure ast.AbstractSyntaxTree, depth int) {

	for _, node := range structure {
		last := node == structure[len(structure)-1]
		const tBoxCharacter string = "  ├ "
		const lBoxCharacter string = "  └ "
		if last {
			// Depth - 1 because a depth of still runs once.
			for range depth {
				fmt.Print(tBoxCharacter)
			}

			fmt.Print(lBoxCharacter)
		} else {
			for range depth {
				fmt.Print(tBoxCharacter)
			}

			fmt.Print(tBoxCharacter)
		}

		fmt.Println(node.Type())

		recursiveVerbosePrint(node.Children(), depth+1)
	}
}
