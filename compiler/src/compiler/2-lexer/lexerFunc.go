package lexer

type LexFunc func(*Lexer) (V2LexNode, LexFunc)
