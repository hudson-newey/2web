package ssg

import "hudson-newey/2web/src/lexer"

var ssgStartToken lexer.LexerToken = []string{"{%"}
var ssgEndToken lexer.LexerToken = []string{"%}"}

// each function will have a keyword that can be used inside ssg contexts
var includeToken lexer.LexerToken = []string{"include"}
var forToken lexer.LexerToken = []string{"for"}
