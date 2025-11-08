package states

const (
	// Text that will be rendered in the browser.
	TextContent LexState = "Text"

	// Text that exists because it didn't match any matchers in the state func.
	SourceText LexState = "SourceText"

	// <div[here]>
	Element LexState = "Element"

	// <!--[here]-->
	MarkupComment LexState = "MarkupComment"

	// /*[here]*/
	ScriptComment LexState = "ScriptComment"

	// <script compiled>[here]</script>
	CompiledScriptSource LexState = "CompiledScriptSource"

	// <script>[here]</script>
	ScriptSource LexState = "ScriptSource"

	// <style>[here]</style>
	StyleSource LexState = "StyleSource"

	// e.g. "[here]" or '[here]'
	String LexState = "String"
)
