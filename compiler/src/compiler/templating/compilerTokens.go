package templating

import "hudson-newey/2web/src/compiler/lexer"

var runtimeScriptStartToken lexer.LexerToken = []string{"<script>"}
var runtimeScriptEndToken lexer.LexerToken = []string{"</script>"}

var styleStartToken lexer.LexerToken = []string{"<style>"}
var styleEndToken lexer.LexerToken = []string{"</style>"}

var variableToken lexer.LexerToken = []string{"$"}

var statementEndToken lexer.LexerToken = []string{";"}
var newLineToken lexer.LexerToken = []string{"\n"}

// e.g. @click="$count = $count + 1"
var eventPrefix lexer.LexerPropPrefix = []string{"@"}

// e.g. *innerText="$count"
var propertyPrefix lexer.LexerPropPrefix = []string{"*"}

// Import suffixes use the "statementEndToken"
// This means that all imports MUST have trailing semi-columns
var importPrefix lexer.LexerToken = []string{"import"}

var lineCommentStartToken lexer.LexerToken = []string{"//"}

var blockCommentStartToken lexer.LexerToken = []string{"/*"}
var blockCommentEndToken lexer.LexerToken = []string{"*/"}

var markupCommentStartToken lexer.LexerToken = []string{"<!--"}
var markupCommentEndToken lexer.LexerToken = []string{"-->"}
