package compiler

import "hudson-newey/2web/src/lexer"

var compilerStartToken lexer.LexerToken = []string{"<script compiled>"}
var compilerEndToken lexer.LexerToken = []string{"</script>"}

var variableToken lexer.LexerToken = []string{"$"}

var statementEndToken lexer.LexerToken = []string{";"}
var variableAssignmentToken lexer.LexerToken = []string{"="}

// e.g. @click="$count = $count + 1"
var eventPrefix lexer.LexerPropPrefix = []string{"@"}

// e.g. *innerText="$count"
var propertyPrefix lexer.LexerPropPrefix = []string{"*"}
