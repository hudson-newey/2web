package lexer

func PrintVerboseLexer(structure []V2LexNode) {
	for _, x := range structure {
		x.PrintDebug()
	}
}
