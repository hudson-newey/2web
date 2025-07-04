package lexerV2State

type State = int

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
