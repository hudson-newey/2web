package parser

import (
	"fmt"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

func PrintVerboseParser(structure nodes.AbstractSyntaxTree) {
	recursiveVerbosePrint(structure, 0)

	// Draw an empty line at the end for better readability and to visually
	// separate this log section from any subsequent logs.
	fmt.Println()
}

func recursiveVerbosePrint(structure nodes.AbstractSyntaxTree, depth int) {

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
