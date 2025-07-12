package preprocessor

import lexer "hudson-newey/2web/src/compiler/2-lexer"

// ssg start and end tokens are exported because we use the same syntax inside
// control flow such as if conditions e.g. {% if <condition> <result> %}
// I have chosen to use the same if start and end tokens for ssg and control
// flow because one of the goals of this project is to make the difference
// between ssg, ssr, isr, and csr an implementation detail that the user
// doesn't have to deal with.
// The compiler should be able to automatically pick the most efficient
// rendering method depending on the circumstances.
var SsgStartToken lexer.LexerSelector = []string{"{%"}
var SsgEndToken lexer.LexerSelector = []string{"%}"}

// each function will have a keyword that can be used inside ssg contexts
var includeToken lexer.LexerSelector = []string{"include"}
