package lexer

// <script compiled>
type CompNode LexNodeType[CompNode]

// $ varName = value
type VarNode LexNodeType[PropNode]

// [propName $varName]
type PropNode LexNodeType[PropNode]
