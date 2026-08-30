package grammar

// Rules that can appear inside of text content.
// These rules are typically "head" nodes that contain other nodes.
var TextRules = []Grammar{
	inlineStyles,

	inlineScripts,
	compiledScripts,

	codeBlock,

	reactiveProperty,
	reactiveEvent,

	elementRef,
	textOutput,

	// Control flow always has to appear last so that exiting a control flow
	// block is always avoided.
	// e.g. otherwise a text output node could match against the end of an if
	// condition because it contains a curly closing brace.
	controlFlowIf,
}
