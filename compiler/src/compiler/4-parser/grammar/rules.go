package grammar

// Rules that can appear inside of text content.
// These rules are typically "head" nodes that contain other nodes.
var TextRules = []grammar{
	inlineStyles,

	inlineScripts,
	compiledScripts,
}
