package lexer

type LexNodeType[T voidNode] any

type LexNode[T LexNodeType[T]] struct {
	Selector string
	Content  string
	Tokens   []string
}
