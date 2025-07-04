package lexerV2Tokens

// For more information, see the possible tokens returned from the parsing
// "Tokenization" step, see the spec:
// https://html.spec.whatwg.org/multipage/parsing.html#tokenization

// We don't use the any helper here because we might want to extend the token
// interface in the future.
type LexNodeToken interface{}
