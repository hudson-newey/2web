package lexer

const (
	// Text that will be rendered in the browser.
	textContent lexState = "Text"

	// Text that exists because it didn't match any matchers in the state func.
	sourceText lexState = "SourceText"

	// <div[here]>
	element lexState = "Element"

	// @here
	controlFlow = "ControlFlow"

	// <!--[here]-->
	markupComment lexState = "MarkupComment"

	// /*[here]*/
	scriptComment lexState = "ScriptComment"

	// <script compiled>[here]</script>
	compiledScriptSource lexState = "CompiledScriptSource"

	// <script>[here]</script>
	scriptSource lexState = "ScriptSource"

	// <style>[here]</style>
	styleSource lexState = "StyleSource"

	// <code>[here]</code>
	codeSource lexState = "CodeSource"

	// e.g. "[here]" or '[here]'
	sourceString lexState = "String"
)
