package lexerV2Tokens

type Attribute struct {
	Name  string
	Value string
}

type StartToken struct {
	SelfClosing bool
	Attributes  []string
}

type EndToken struct {
	SelfClosing bool
	Attributes  []string
}
