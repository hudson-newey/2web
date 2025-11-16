package grammar

import (
	lexerTokens "hudson-newey/2web/src/compiler/2-lexer/tokens"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var scriptImport = grammar{
	Def: newDefinition(
		lexerTokens.KeywordImport,
		lexerTokens.TextContent,
		lexerTokens.KeywordFrom,
		anyQuote,
		lexerTokens.TextContent,
		anyQuote,
		lexerTokens.Semicolon,
	),
	Constructor: wrapConstructor(nodes.NewScriptImportNode),
	ChildDefs:   []grammar{},
}
