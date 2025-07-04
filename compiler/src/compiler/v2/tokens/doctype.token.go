package lexerV2Tokens

type DoctypeToken struct {
	Name             string
	PublicIdentifier string
	SystemIdentifier string
	ForceQuirks      bool
}
