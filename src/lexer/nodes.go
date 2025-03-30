package lexer

// <script compiled>
type CompNode LexNodeType[CompNode]

// $ varName = value
type VarNode LexNodeType[PropNode]

// [propName $varName]
type PropNode LexNodeType[PropNode]

// @eventName $varName = value@
type EventNode LexNodeType[EventNode]

// {% ssgModule arguments %}
type SsgNode LexNodeType[SsgNode]
