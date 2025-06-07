package lexer

// <script compiled>
type CompNode LexNodeType[CompNode]

// $ varName = value
type VarNode LexNodeType[PropNode]

// [propName $varName]
type PropNode LexNodeType[PropNode]

// @eventName $varName = value@
type EventNode LexNodeType[EventNode]

// {{ $variable }}
type TextNode LexNodeType[TextNode]

// {% ssgModule arguments %}
type SsgNode LexNodeType[SsgNode]

// #elementRef
type RefNode LexNodeType[RefNode]

// import Component from "components/footer.component.html"
type ImportNode LexNode[ImportNode]

// A comment starting with "//"
type LineCommentNode LexNode[LineCommentNode]

// A comment with start /* and end */ delimiters
type BlockCommentNode LexNode[BlockCommentNode]

// A comment with <!-- and end --> delimiters
type MarkupCommentNode LexNode[MarkupCommentNode]

type voidNode LexNodeType[struct{}]
