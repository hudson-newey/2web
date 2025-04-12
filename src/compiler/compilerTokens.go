package compiler

import "hudson-newey/2web/src/lexer"

var attributeEndToken lexer.LexerToken = []string{"*", ">", " ", "\n", "/"}
var statementEndToken lexer.LexerToken = []string{";"}
var variableAssignmentToken lexer.LexerToken = []string{"="}

var compilerStartToken lexer.LexerToken = []string{"<script compiled>"}
var compilerEndToken lexer.LexerToken = []string{"</script>"}

var variableToken lexer.LexerToken = []string{"$"}

var reactiveStartToken lexer.LexerToken = []string{"["}
var reactiveEndToken lexer.LexerToken = []string{"]"}

var mustacheStartToken lexer.LexerToken = []string{"{{ "}
var mustacheEndToken lexer.LexerToken = []string{" }}"}

var eventStartToken lexer.LexerToken = []string{"@"}
