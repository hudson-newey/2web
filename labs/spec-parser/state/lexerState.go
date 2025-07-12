package lexerV2State

import (
	lexerV2Errors "hudson-newey/2web/src/compiler/v2/errors"
	lexerV2Tokens "hudson-newey/2web/src/compiler/v2/tokens"
)

type State = int
type NestingLevel = int
type ParsePauseState = bool
type ReturnState = State

type NextInputCharacter = rune

type HandlerReturn = func(char NextInputCharacter, nesting NestingLevel) (State, *ReturnState, lexerV2Tokens.LexNodeToken, *lexerV2Errors.ParseError)

// https://html.spec.whatwg.org/multipage/parsing.html#tokenization
const (
	Data State = iota
	RCDATA

	RAWTEXT
	ScriptData
	PLAINTEXT
	TagOpen
	EndTagOpen
	TagName

	RCDATALessThanSign
	RCDATAEndTagOpen
	RCDATAEndTagName

	RAWTEXTLessThanSign
	RAWTEXTEndTagOpen
	RAWTEXTEndTagName

	ScriptDataLessThanSign
	ScriptDataEndTagOpen
	ScriptDataEndTagName
	ScriptDataEscapeStart
	ScriptDataStartDashState
	ScriptDataEscapedState
	ScriptDataEscapedDash
	ScriptDataEscapedDashDash
	ScriptDataEscapedLessThanSign
	ScriptDataEscapedEndTagOpen
	ScriptDataEscapedEndTagName

	ScriptDataDoubleEscapeStart
	ScriptDataDoubleEscaped
	ScriptDataDoubleEscapedDash
	ScriptDataDoubleEscapedDashDash
	ScriptDataDoubleEscapedLessThan
	ScriptDataDoubleEscapedEnd

	BeforeAttributeName
	AttributeName
	AfterAttributeName
	BeforeAttributeValue
	AttributeValueDoubleQuoted
	AttributeValueSingleQuoted
	AttributeValueUnquotedQuoted
	AfterAttributeValueQuoted
	SelfClosingStartTag

	BogusComment
	MarkupDeclaration

	CommentStart
	CommentStartDash
	Comment
	CommentLessThanSign
	CommentLessThanSignBang
	CommentLessThanSignBangDash
	CommentLessThanSignBangDashDash
	CommentEndDash
	CommentEnd
	CommentEndBang

	DOCTYPE
	BeforeDoctypeName
	DoctypeName
	AfterDoctypeName
	AfterDoctypePublicKeyword
	BeforeDoctypePublicIdentifier
	DoctypePublicIdentifierDoubleQuoted
	DoctypePublicIdentifierSingleQuoted
	AfterDoctypePublicIdentifier
	BetweenDoctypePublicAndSystemIdentifier
	AfterDoctypeSystemKeyword
	BeforeDoctypeIdentifierKeyword
	DoctypeSystemIdentifierDoubleQuoted
	DoctypeSystemIdentifierSingleQuoted
	AfterDoctypeSystemIdentifier
	BogusDoctype

	// TODO: CDATA

	CharacterReference
	NamedCharacterReference
	AmbiguousAmpersand
	NumericCharacterReference
	HexadecimalCharacterReferenceStart
	DecimalCharacterReferenceStart
	HexadecimalCharacterReference
	DecimalCharacterReference
	NumericCharacterReferenceEnd
)
