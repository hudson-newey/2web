package lexer

import (
	"os"
	"text/tabwriter"
)

func PrintVerboseLexer(structure []V2LexNode) {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	defer tabWriter.Flush()

	tableHeader := "Row:Col\tState\tToken\tContent\n"
	_, err := tabWriter.Write([]byte(tableHeader))
	if err != nil {
		panic(err)
	}

	tabWriter.Write([]byte("----\t-----\t-----\t-------\n"))

	for _, x := range structure {
		bytesToWrite := []byte(x.PrintDebug())

		_, err := tabWriter.Write(bytesToWrite)
		if err != nil {
			panic(err)
		}
	}
}
