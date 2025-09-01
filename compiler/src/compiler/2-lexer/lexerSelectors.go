package lexer

// Represents a lexical token that can be used in the lexer
//
// The LexerSelector is a string array because a lexer token may have one or more
// variations. For instance, the terminating token in attribute selector can
// either be an opening carrot, space, a new line, or a slash
//
// E.g. This is all valid attribute syntax
//
//	<button @click="$count++">+1</button>
//
//	<button title="increment" @click="$count++">+1</button>
//
//	<button
//	 title="increment"
//	 @click="$count++"
//	>
//	 +1
//	</button>
//
//	<img @load="onImageLoad()"/>
//
// _Note: this last example of a slash being a valid terminator is due to self_
// closing tags.
type LexerSelector []string
