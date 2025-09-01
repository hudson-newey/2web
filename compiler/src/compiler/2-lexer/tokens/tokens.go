package lexerTokens

type LexToken int

const (
	// EOF
	EOF LexToken = iota

	// Used in doctypes
	// !
	Exclamation

	// <
	LessAngle

	// >
	GreaterAngle

	// /
	Slash

	// =
	Equals

	// '
	QuoteSingle

	// "
	QuoteDouble

	// <!--
	CommentStart

	// -->
	CommentEnd

	// All tokens below are custom tokens to 2web

	// @
	AtSymbol

	// *
	Star

	// #
	Hash

	Text

	// JS Tokens

	// import
	KeywordImport

	// $
	DollarSign

	// ;
	SemiColon

	// //
	LineCommentStart

	// /*
	BlockCommentStart

	// */
	BlockCommentEnd
)
