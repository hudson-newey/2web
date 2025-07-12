package templating

import lexer "hudson-newey/2web/src/compiler/2-lexer"

var runtimeScriptStartToken = lexer.LexerSelector{"<script>"}
var runtimeScriptEndToken = lexer.LexerSelector{"</script>"}

var styleStartToken = lexer.LexerSelector{"<style>"}
var styleEndToken = lexer.LexerSelector{"</style>"}

var variableToken = lexer.LexerSelector{"$"}

var statementEndToken = lexer.LexerSelector{";"}
var newLineToken = lexer.LexerSelector{"\n"}

// e.g. @click="$count = $count + 1"
var eventPrefix = lexer.LexerSelector{"@"}

// e.g. *innerText="$count"
var propertyPrefix = lexer.LexerSelector{"*"}

// Import suffixes use the "statementEndToken"
// This means that all imports MUST have trailing semi-columns
var importPrefix = lexer.LexerSelector{"import"}

var lineCommentStartToken = lexer.LexerSelector{"//"}

var blockCommentStartToken = lexer.LexerSelector{"/*"}
var blockCommentEndToken = lexer.LexerSelector{"*/"}

var markupCommentStartToken = lexer.LexerSelector{"<!--"}
var markupCommentEndToken = lexer.LexerSelector{"-->"}
