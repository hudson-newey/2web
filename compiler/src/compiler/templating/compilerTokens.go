package templating

import "hudson-newey/2web/src/compiler/lexer"

var compilerScriptStartToken lexer.LexerToken = []string{"<script compiled>"}
var compilerScriptEndToken lexer.LexerToken = []string{"</script>"}

var runtimeScriptStartToken lexer.LexerToken = []string{"<script>"}
var runtimeScriptEndToken lexer.LexerToken = []string{"</script>"}

var styleStartToken lexer.LexerToken = []string{"<style>"}
var styleEndToken lexer.LexerToken = []string{"</style>"}

var variableToken lexer.LexerToken = []string{"$"}

var statementEndToken lexer.LexerToken = []string{";"}

// e.g. @click="$count = $count + 1"
var eventPrefix lexer.LexerPropPrefix = []string{"@"}

// e.g. *innerText="$count"
var propertyPrefix lexer.LexerPropPrefix = []string{"*"}

// Import suffixes use the "statementEndToken"
// This means that all imports MUST have trailing semi-columns
var importPrefix lexer.LexerToken = []string{"import"}
