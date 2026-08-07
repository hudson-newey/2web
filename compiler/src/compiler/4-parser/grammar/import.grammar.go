package grammar

import (
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/nodes"
)

var scriptImport = Grammar{
	Def: newDefinition(
		lexeme.KeywordImport,
		lexeme.CompiledScriptSource,
		lexeme.KeywordFrom,
		lexeme.CompiledScriptSource,
		lexeme.Semicolon,
	),
	Constructor: wrapConstructor(nodes.NewscriptImportNode),
	ChildDefs:   []Grammar{},
}
