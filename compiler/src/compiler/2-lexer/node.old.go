package lexer

// TODO: I am currently in the process of phasing out this lexer representation
// because its usage is poorly architectured
type LexNode[T LexNodeType[T]] struct {
	Selector string
	Content  string
	Tokens   []string
}
