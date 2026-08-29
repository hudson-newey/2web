package lexeme

const (
	// EOF
	EOF Lexeme = "EOF"

	// I preserve newlines as their own token because I think that having the
	// output be as close as possible to the input is a good thing for development
	// environments, because it makes it easier to debug.
	// \n
	Newline Lexeme = "Newline"

	// !DOCTYPE
	Doctype Lexeme = "Doctype"

	// <
	LessAngle Lexeme = "LessAngle"

	// >
	GreaterAngle Lexeme = "GreaterAngle"

	// /
	Slash Lexeme = "Slash"

	// =
	Equals Lexeme = "Equals"

	// '
	QuoteSingle Lexeme = "QuoteSingle"

	// "
	QuoteDouble Lexeme = "QuoteDouble"

	// `
	Backtick Lexeme = "Backtick"

	// Space character (" ")
	Space Lexeme = "Space"

	// <!--
	MarkupCommentStart Lexeme = "MarkupCommentStart"

	// -->
	MarkupCommentEnd Lexeme = "MarkupCommentEnd"

	// All tokens below are custom tokens to 2web

	// @
	AtSymbol Lexeme = "AtSymbol"

	// *
	Star Lexeme = "Star"

	// !
	Exclamation Lexeme = "Exclamation"

	// #
	Hash Lexeme = "Hash"

	// {
	CurlyOpen Lexeme = "CurlyOpen"

	// }
	CurlyClosed Lexeme = "CurlyOpen"

	// Backslash used for escaping
	// (e.g. \> in HTML so that you don't have to use the ugly &gt;)
	// \
	Escape Lexeme = "Escape"

	// <style>
	StyleStartTag Lexeme = "StyleStartTag"

	// </style>
	StyleEndTag Lexeme = "StyleEndTag"

	// <script compiled>
	CompiledScriptStartTag Lexeme = "CompiledScriptStartTag"

	// <script>
	ScriptStartTag Lexeme = "ScriptStartTag"

	// </script>
	ScriptEndTag Lexeme = "ScriptEndTag"

	//? JS Tokens

	// import
	KeywordImport Lexeme = "KeywordImport"

	// from
	// This is typically used together with the import keyword.
	KeywordFrom Lexeme = "KeywordFrom"

	// $
	DollarSign Lexeme = "DollarSign"

	// ;
	Semicolon Lexeme = "Semicolon"

	// //
	ScriptLineCommentStart Lexeme = "ScriptLineCommentStart"

	// newline
	ScriptLineCommentEnd Lexeme = "ScriptLineCommentEnd"

	// /*
	ScriptBlockCommentStart Lexeme = "ScriptBlockCommentStart"

	// */
	BlockCommentEnd Lexeme = "BlockCommentEnd"

	//? General

	// TextContent content that is displayed on the screen
	TextContent Lexeme = "TextContent"

	// Inside a <script compiled> tag or in an external .js file
	CompiledScriptSource Lexeme = "CompiledScriptSource"

	// Inside a <script> tag or in an external .js file
	ScriptSource Lexeme = "ScriptSource"

	// Inside a <style> tag or in an external .css file
	StyleSource Lexeme = "StyleSource"

	// Code blocks are distinct from regular HTML elements in 2web because they
	// are automatically escaped.
	// I have made this distinction at the lexer level to avoid ambiguity.

	// <code>
	CodeStartTag Lexeme = "CodeStartTag"

	// Inside a <code> tag
	CodeSource Lexeme = "CodeSource"

	// </code>
	CodeEndTag Lexeme = "CodeEndTag"
)
