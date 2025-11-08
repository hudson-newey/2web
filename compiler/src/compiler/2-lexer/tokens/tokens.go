package lexerTokens

const (
	// EOF
	EOF LexToken = "EOF"

	// I preserve newlines as their own token because I think that having the
	// output be as close as possible to the input is a good thing for development
	// environments, because it makes it easier to debug.
	// \n
	Newline LexToken = "Newline"

	// !DOCTYPE
	Doctype LexToken = "Doctype"

	// <
	LessAngle LexToken = "LessAngle"

	// >
	GreaterAngle LexToken = "GreaterAngle"

	// /
	Slash LexToken = "Slash"

	// =
	Equals LexToken = "Equals"

	// '
	QuoteSingle LexToken = "QuoteSingle"

	// "
	QuoteDouble LexToken = "QuoteDouble"

	// `
	Backtick LexToken = "Backtick"

	// Space character (" ")
	Space LexToken = "Space"

	// <!--
	MarkupCommentStart LexToken = "MarkupCommentStart"

	// -->
	MarkupCommentEnd LexToken = "MarkupCommentEnd"

	// All tokens below are custom tokens to 2web

	// @
	AtSymbol LexToken = "AtSymbol"

	// *
	Star LexToken = "Star"

	// !
	Exclamation LexToken = "Exclamation"

	// #
	Hash LexToken = "Hash"

	// Backslash used for escaping
	// (e.g. \> in HTML so that you don't have to use the ugly &gt;)
	// \
	Escape LexToken = "Escape"

	//? JS Tokens

	// import
	KeywordImport LexToken = "KeywordImport"

	// $
	DollarSign LexToken = "DollarSign"

	// ;
	SemiColon LexToken = "SemiColon"

	// <style>
	StyleStartTag LexToken = "StyleStartTag"

	// </style>
	StyleEndTag LexToken = "StyleEndTag"

	// <script>
	ScriptStartTag LexToken = "ScriptStartTag"

	// </script>
	ScriptEndTag LexToken = "ScriptEndTag"

	// //
	ScriptLineCommentStart LexToken = "ScriptLineCommentStart"

	// //
	ScriptLineCommentEnd LexToken = "ScriptLineCommentEnd"

	// /*
	ScriptBlockCommentStart LexToken = "ScriptBlockCommentStart"

	// */
	BlockCommentEnd LexToken = "BlockCommentEnd"

	//? General

	// TextContent content that is displayed on the screen
	TextContent LexToken = "Text"

	// Inside a <script> tag or in an external .js file
	ScriptSource LexToken = "ScriptSource"

	// Inside a <style> tag or in an external .css file
	StyleSource LexToken = "StyleSource"
)
