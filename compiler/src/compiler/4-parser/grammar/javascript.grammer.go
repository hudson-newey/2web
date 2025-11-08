package grammar

import lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"

var ReactiveVariable []lexerTokens.LexToken = []lexerTokens.LexToken{
	lexerTokens.DollarSign,
	lexerTokens.Space,
	lexerTokens.TextContent, // variable name
	lexerTokens.SemiColon,
}
