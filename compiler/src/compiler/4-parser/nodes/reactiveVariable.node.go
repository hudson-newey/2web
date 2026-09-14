package nodes

import (
	"fmt"
	"hudson-newey/2web/src/cli"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/document/documentErrors"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
	"hudson-newey/2web/src/models"
	"strings"

	"github.com/hudson-newey/2web/_shared/lists"
	"github.com/hudson-newey/2web/_shared/logger"
)

func NewreactiveVariableNode(lexNodes []*lexer.V2LexNode, context *ParseContext) Node {
	variableName, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 1)
	if err != nil {
		return context.DegradedNode(
			"reactive variable declaration is missing a variable name. Reactive variables are declared with '$name = value;'",
			lexNodes,
		)
	}

	initialValue, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 2)
	if err != nil {
		return context.DegradedNode(
			"reactive variable declaration is missing an initial value. Reactive variables are declared with '$name = value;'",
			lexNodes,
		)
	}

	variableNameContent := strings.TrimSpace(variableName.Content)
	initialValueContent := strings.TrimSpace(initialValue.Content)

	if initialValueContent == "" {
		return context.DegradedNode(
			"reactive variable '$"+variableNameContent+"' is missing an initial value",
			lexNodes,
		)
	}

	return &reactiveVariableNode{
		variableName: variableNameContent,
		initialValue: initialValueContent,
		// The declaration position is kept so that errors about the variable
		// (e.g. that it is unused) can be reported against the declaration.
		position: positionOf(lexNodes),
	}
}

type reactiveVariableNode struct {
	variableName string
	initialValue string
	children     AbstractSyntaxTree

	// The position of the declaration in the source file.
	position lexer.Position
}

func (m *reactiveVariableNode) Type() string {
	return "reactiveVariableNode"
}

func (m *reactiveVariableNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *reactiveVariableNode) MarkupContent() string {
	return ""
}

func (m *reactiveVariableNode) Content(page *page.Page, index *ReactiveIndex) NodeContent {
	if !cli.GetArgs().NoReactivity {
		// TODO: This should be a non-mutative operation
		m.compileReactivity(page, index)
	}

	return NodeContent{
		HtmlContent:      page.Html,
		JsContent:        javascript.NewJsFile(),
		CssContent:       css.NewCssFile(),
		TwoScriptContent: twoscript.NewTwoScriptFile(),
	}
}

func (m *reactiveVariableNode) AddChild(child Node) {
	m.children = append(m.children, child)
}

func (m *reactiveVariableNode) RemoveChild(child Node) {
	for i, c := range m.children {
		if c == child {
			m.children = append(m.children[:i], m.children[i+1:]...)
			return
		}
	}
}

// Selector returns the variable's selector (its name prefixed with a dollar
// sign).
func (m *reactiveVariableNode) Selector() string {
	return m.selector()
}

// When using this variable in code, reactive variables are easily visually
// separated from normal JavaScript variables using the dollar sign prefix.
func (m *reactiveVariableNode) selector() string {
	if strings.HasPrefix(m.variableName, "$") {
		logger.PrintWarning(
			fmt.Sprintf("Reactive variable declaration '%s' has double $$ prefix.", m.variableName),
		)
	}
	return fmt.Sprintf("$%s", m.variableName)
}

func (m *reactiveVariableNode) dependentProps(index *ReactiveIndex) []*reactivePropertyNode {
	return index.FindDependentProperties(m)
}

func (m *reactiveVariableNode) dependentEvents(index *ReactiveIndex) []*reactiveEventNode {
	events := index.FindDependentEvents(m)

	// Event reducers that directly call an imported server function don't
	// assign to any variable. They are compiled into standalone listeners by
	// CompileServerCalls and are never wired here (wiring them here would
	// attach the call to the assignment machinery of a variable it merely
	// reads).
	return lists.Filter(events, func(e *reactiveEventNode) bool {
		return e.assignmentSink != ""
	})
}

type reactivityLevel int

/*
Reactive types progressively get less performant as you go down this list.

"Static" reactive variables are not really reactive at all, and we can
inline them directly at compile time.

"StaticProperty"
Variables that require initial bootstrapping.
If a variable is a StaticProperty, it means that the variable does not change
after the initial render, but it requires a <script> tag to modify the DOM
on initial render.

e.g.
```html
<script compiled>
$ message = "Hello!";
</script>

<h1 *textContent="$message"></h1>
```

As you can see from the example, the variable is not really reactive, but it
does require a <script> tag to modify the DOM on initial render.

The compiled code will look something like this:

```html
<h1 id="_0">Hello!</h1>
<script>

	document.addEventListener("DOMContentLoaded", () => {
	    document.getElementById("_0").textContent = "Hello!";
	});

</script>
```

"Assignment" reactive variables are reactive, but do not require a runtime
variable to keep track of state.
e.g.
```html
<script compiled>
$ message = "Hello!";
</script>

<p>{{ $message }}</p>

<button @click="$message = 'World'">Change message</button>
```

In this example, we don't need to keep track of the $message state, because it
it's mutation does not depend on its previous value.
We can just directly replace the <p> tags content with the "World" string when
the button is clicked.

The compiled code will look something like this:
```html
<p id="_0">Hello!</p>
<button onclick="0()">Change message</button>
<script>

	function 0() {
	    document.getElementById("_0").innerHTML = "World";
	}

</script>
```

As you can see, there is no additional runtime variable to keep track of the
"$message" state.

"Reactive" reactive variables are the least performant, because they require a
runtime variable to keep track of state.
Note: This is not a signal, or any other reactive state such as a Proxy object.
This is just a simple "let" variable that we can mutate and read from.

This is much more performant than a signal, but should be avoided if possible.

We only need "Reactive" variable types when updating a variable depends on its
previous value.

e.g.
```html
<script compiled>
$ count = 0;
</script>

<p>{{ $count }}</p>
<button @click="$count++">Increment</button>

In this example, out compiled code will look **something** (not exact) like this

```html
<p id="_0">0</p>

<button @click="0()">Increment</button>

<script>

	let count = 0;

	function 0() {
	    document.getElementById("_0").innerHTML = count++;
	}

</script>
```
*/
const (
	// never used throughout the app. A noop.
	unused reactivityLevel = iota

	// does not require any JavaScript. Can be inlined at compile time.
	static

	// Requires JavaScript to modify the DOM on initial render.
	staticProperty

	// Requires JavaScript to attach an event listener to the DOM and modify
	// the DOM on event.
	assignment

	// Requires JavaScript to attach an event listener, keep track of state,
	// and modify the DOM on event.
	reactive

	// Computed from other reactive variables. Its value is re-evaluated by
	// the update cascade of the variables it is computed from.
	computed
)

// TODO: this should probably cache the type for faster compile times
func (m *reactiveVariableNode) reactivityLevel(index *ReactiveIndex) reactivityLevel {
	// A variable that other variables are computed from needs a runtime
	// representation so that the computed variables can be re-evaluated when
	// it changes.
	if index.HasDerivedDependents(m) {
		return reactive
	}

	if index.IsDerived(m) {
		return computed
	}

	// A variable that is assigned inside a compiled script function needs a
	// runtime representation and an update function: the function body's
	// assignments are rewritten into "update the runtime variable and run its
	// update cascade" statements (see twoScriptNode).
	if index.IsAssignedByFunction(m) {
		return reactive
	}

	events := m.dependentEvents(index)
	for _, e := range events {
		// If the assignment expression uses the same variable that it's
		// assigning to, we need to have a runtime variable to track state.
		// e.g. think of a counting number
		if strings.Contains(e.assignmentExpr, m.selector()) {
			return reactive
		}
	}

	// If we need to respond to user input but the responses don't depend on
	// previous state, we can keep updates completely stateless.
	// e.g. think of a close button on a dialog. We can just set closed to false
	// without needing to know if the dialog is open.
	if len(events) > 0 {
		return assignment
	}

	// All reactive properties that require an initial runtime assignment, but
	// don't ever update after first page load.
	// e.g. think of a date/time that can't be evaluated at runtime.
	props := m.dependentProps(index)
	hasNonOptimizableProps := lists.Some(props, func(x *reactivePropertyNode) bool {
		return !x.canCompilerInline()
	})
	// If even one of the properties cannot be inlined, then we have to treat it
	// as a static property.
	// TODO: If we start processing static reducers on the prop level, we can
	// probably do per-prop optimizations.
	if hasNonOptimizableProps {
		return staticProperty
	}

	if len(props) > 0 {
		return static
	}

	return unused
}

func (m *reactiveVariableNode) compileReactivity(pageModel *page.Page, index *ReactiveIndex) {
	// short circuit fast if not used
	if index.IsUnused(m) {
		errMsg := fmt.Sprintf("Unused variable: %s", m.selector())
		err := models.NewError(errMsg, pageModel.InputPath, m.position)
		documentErrors.AddErrors(&err)
		return
	}

	// A variable that is computed from other reactive variables is not wired
	// here: it is either recomputed by the update cascade of the runtime
	// variable it depends on, or (when all of its dependencies are static) it
	// is evaluated once in a startup bootstrap script.
	if index.IsDerived(m) {
		// A computed variable that is also assigned by a compiled script
		// function must keep its runtime representation (the bootstrap below
		// evaluates the value once, which would silently drop the function's
		// updates).
		if !index.HasRuntimeDependency(m) && !index.IsAssignedByFunction(m) {
			m.compileComputedBootstrap(pageModel, index)
			return
		}

		if !index.IsRuntime(m) {
			m.compileComputedChunk(pageModel, index)
			return
		}

		// The variable is computed from other variables AND directly assigned
		// by its own events. It compiles like any other runtime variable, with
		// its computation dependencies resolved into the shared reactive
		// runtime scope.
	}

	// Ideally, slower reactive types would only target properties and events
	// that are effected.
	// Therefore, the fully "models.Reactive" variable reactivity class is a
	// superset of "models.Assignment" because some references to the variable
	// might not be fully reactive.
	//
	// This also means that each reactivity type should also make their own
	// element selectors so that subset of elements that abide by the reactivity
	// class are updated.
	//
	// TODO: I might be able to combine selectors for the same element that has
	// different property targets.
	if reactivityLevel := m.reactivityLevel(index); reactivityLevel >= reactive {
		m.compileReactiveVar(pageModel, index)
	} else if reactivityLevel >= assignment {
		// TODO: explore if reactive and assignment reactivity are mutually
		// exclusive for variables, events, or props
		m.compileAssignmentVar(pageModel, index)
	}

	// static props differ from truly static variables because static props
	// need runtime code to set the initial value of the prop
	//
	// e.g. <my-custom-element *value="$value"></my-custom-element>
	//
	// Most elements that have writable properties, also have an associated
	// attribute. So static properties only really apply to poorly designed
	// custom elements (e.g. web components).
	//
	// If the star is removed from this attribute, then it will become static
	// content that can be evaluated at runtime.
	// It is therefore recommended to remove the star from attributes if you
	// do not need to
	//
	// e.g. <input type="range" value="$value"></input>
	if reactivityLevel := m.reactivityLevel(index); reactivityLevel >= staticProperty {
		m.compileStaticPropVar(pageModel, index)
	}

	if reactivityLevel := m.reactivityLevel(index); reactivityLevel == static {
		m.compileStatic(pageModel, index)
	}
}

func (m *reactiveVariableNode) compileReactiveVar(
	pageModel *page.Page,
	index *ReactiveIndex,
) {
	events := m.dependentEvents(index)
	props := m.dependentProps(index)

	jsNewValueVar := javascript.ValueVar

	domMutator := ""
	for _, p := range props {
		domMutator = domMutator + fmt.Sprintf(
			`document.querySelectorAll("[%s]").forEach((__2_element_ref_mod) => __2_element_ref_mod["%s"] = %s);`,
			p.selector(pageModel), p.propName, p.propAssignment(p.propValue(jsNewValueVar)),
		)
	}

	// The runtime variable name was pre-allocated before the reactivity pass
	// (see templating.Compile) so that computed variables and compiled script
	// functions can reference it regardless of compilation order.
	variableName := index.RuntimeVariableName(m.selector())

	// The update function's name is pre-allocated with it, so that compiled
	// script functions can call the update cascade of the variables they
	// assign to.
	handlerFuncName := index.HandlerName(m.selector())
	if handlerFuncName == "" {
		handlerFuncName = pageModel.Ids.CreateFunctionName()
	}

	// Variables that are computed from this one are re-evaluated by this
	// variable's update cascade: whenever an event changes this variable, every
	// computed variable that (transitively) depends on it is re-evaluated (in
	// dependency order) and its bound properties are updated.
	domMutator = domMutator + m.compileComputedCascade(pageModel, index)

	domMutator = fmt.Sprintf(`
	 	let %s = %s;
		function %s(%s) { %s }
	`,
		variableName, m.initialValue,
		handlerFuncName, jsNewValueVar, domMutator,
	)

	eventListeners := ""
	pageContent := pageModel.Html.Content
	for _, e := range events {
		// The assignment expression (rather than the raw reducer) is used so
		// that the increment/decrement shorthand expands into a correct
		// assignment (e.g. "$count++" becomes "__2_var = __2_var + 1" rather
		// than the no-op "__2_var = __2_var++").
		reactiveReducer := strings.ReplaceAll(e.assignmentExpr, m.selector(), variableName)

		eventDomSelector := pageModel.Ids.CreateElementName()
		pageContent = strings.ReplaceAll(pageContent, e.selector(), eventDomSelector)
		eventListeners = eventListeners + fmt.Sprintf(`
			document.querySelector("[%s]").addEventListener("%s", () => {
				%s = %s;
				%s(%s);
			});
		`,
			eventDomSelector, e.eventName,
			variableName, reactiveReducer,
			handlerFuncName, variableName,
		)
	}

	handlerContent := fmt.Sprintf("%s\n%s", domMutator, eventListeners)
	pageModel.SetContent(pageContent)

	pageModel.AppendReactiveRuntime(page.ReactiveRuntimeChunk{
		Variable:     m.selector(),
		Dependencies: index.DependenciesOf(m),
		Content:      handlerContent,
	})
}

func (m *reactiveVariableNode) compileAssignmentVar(
	pageModel *page.Page,
	index *ReactiveIndex,
) {
	events := m.dependentEvents(index)
	props := m.dependentProps(index)

	jsNewValueVar := javascript.ValueVar

	domMutator := ""
	for _, p := range props {
		domMutator = domMutator + fmt.Sprintf(
			`document.querySelectorAll("[%s]").forEach((__2_element_ref_mod) => __2_element_ref_mod["%s"] = %s);`,
			p.selector(pageModel), p.propName, p.propAssignment(p.propValue(jsNewValueVar)),
		)
	}

	handlerFuncName := pageModel.Ids.CreateFunctionName()
	domMutator = fmt.Sprintf(
		`function %s(%s) { %s }`,
		handlerFuncName, jsNewValueVar, domMutator,
	)

	eventListeners := ""
	pageContent := pageModel.Html.Content
	for _, e := range events {
		eventDomSelector := pageModel.Ids.CreateElementName()
		pageContent = strings.ReplaceAll(pageContent, e.selector(), eventDomSelector)
		eventListeners = eventListeners + fmt.Sprintf(
			`document.querySelector("[%s]").addEventListener("%s", () => %s(%s));`,
			eventDomSelector, e.eventName, handlerFuncName, e.assignmentExpr,
		)
	}

	handlerContent := fmt.Sprintf("%s\n%s", domMutator, eventListeners)
	handlerScript := javascript.FromGeneratedContent(handlerContent)
	pageModel.SetContent(pageContent)
	pageModel.AddScript(handlerScript)
}

func (m *reactiveVariableNode) compileStaticPropVar(
	pageModel *page.Page,
	index *ReactiveIndex,
) {
	props := m.dependentProps(index)

	reducerContent := ""
	for _, p := range props {
		reducerContent = reducerContent + fmt.Sprintf(
			`document.querySelector("[%s]")["%s"] = %s;`,
			p.selector(pageModel), p.propName, p.propAssignment(p.propValue(m.initialValue)),
		)
	}

	// p.selector() replaces the compile time selector in the page's content, so
	// the page model must not be overwritten with a content snapshot taken
	// before the replacements. Overwriting would leave the compiler selectors
	// in the emitted HTML while the emitted JavaScript queries the runtime
	// selectors, leaving the property permanently unwired.
	reducerScript := javascript.FromGeneratedContent(reducerContent)
	pageModel.AddScript(reducerScript)
}

func (m *reactiveVariableNode) compileStatic(
	pageModel *page.Page,
	index *ReactiveIndex,
) {
	props := m.dependentProps(index)
	for _, p := range props {
		pageModel.SetContent(
			strings.ReplaceAll(pageModel.Html.Content, p.selector(pageModel), m.initialValue),
		)
	}
}

// compileComputedCascade emits the re-evaluation statements for every
// variable that is (transitively) computed from this variable.
//
// The statements are appended to this variable's update function: when an
// event changes this variable, each computed variable is re-evaluated (in
// dependency order, so a computed variable that another computed variable
// depends on is re-evaluated first) and its bound properties are updated.
func (m *reactiveVariableNode) compileComputedCascade(
	pageModel *page.Page,
	index *ReactiveIndex,
) string {
	computedVariables := index.ComputedDependents(m)

	cascade := ""
	for _, computed := range computedVariables {
		runtimeName := index.RuntimeVariableName(computed.selector())

		// A computed dependent without a runtime representation took the
		// bootstrap path (its dependencies are all static, so its value can
		// never change). There is nothing to re-evaluate, and emitting its
		// update would reference an empty variable name.
		if runtimeName == "" {
			continue
		}

		expression, ok := index.ResolveExpression(computed)
		if !ok {
			continue
		}

		cascade = cascade + fmt.Sprintf("\t\t%s = %s;\n", runtimeName, expression)
		cascade = cascade + m.computedPropertyUpdates(pageModel, index, computed, runtimeName)
	}

	return cascade
}

// computedPropertyUpdates emits the DOM property assignments that bind a
// computed variable's value to the elements it is rendered into.
func (m *reactiveVariableNode) computedPropertyUpdates(
	pageModel *page.Page,
	index *ReactiveIndex,
	computed *reactiveVariableNode,
	runtimeName string,
) string {
	updates := ""
	for _, p := range index.FindDependentProperties(computed) {
		updates = updates + fmt.Sprintf(
			`document.querySelectorAll("[%s]").forEach((__2_element_ref_mod) => __2_element_ref_mod["%s"] = %s);`+"\n",
			p.selector(pageModel), p.propName, p.propAssignment(runtimeName),
		)
	}

	return updates
}

// compileComputedChunk wires a computed variable into the reactive runtime
// scope of the variable it is computed from.
//
// The computed variable's runtime declaration, initial DOM assignment, and
// update function are emitted into the shared reactive runtime so that the
// dependency's update cascade can re-evaluate it.
func (m *reactiveVariableNode) compileComputedChunk(
	pageModel *page.Page,
	index *ReactiveIndex,
) {
	expression, ok := index.ResolveExpression(m)
	if !ok {
		errorModel := models.NewError(
			fmt.Sprintf("computed variable '%s' is part of a cyclic dependency", m.selector()),
			pageModel.InputPath,
			m.position,
		)

		documentErrors.AddErrors(&errorModel)
		return
	}

	// The runtime variable name was pre-allocated for this variable.
	runtimeName := index.RuntimeVariableName(m.selector())

	handlerFuncName := pageModel.Ids.CreateFunctionName()

	initialAssignments := m.computedPropertyUpdates(pageModel, index, m, runtimeName)
	mutatorBody := m.computedPropertyUpdates(pageModel, index, m, javascript.ValueVar)

	content := fmt.Sprintf(
		"let %s = %s;\n%s\nfunction %s(%s) { %s }\n",
		runtimeName, expression,
		initialAssignments,
		handlerFuncName, javascript.ValueVar, mutatorBody,
	)

	pageModel.AppendReactiveRuntime(page.ReactiveRuntimeChunk{
		Variable:     m.selector(),
		Dependencies: index.DependenciesOf(m),
		Content:      content,
	})
}

// compileComputedBootstrap evaluates a computed variable whose dependencies
// are all static (none of them can ever change at runtime).
//
// The computed expression is resolved with every dependency inlined as its
// initial value and evaluated once in a startup script. No runtime state is
// needed because the value can never change.
func (m *reactiveVariableNode) compileComputedBootstrap(
	pageModel *page.Page,
	index *ReactiveIndex,
) {
	expression, ok := index.ResolveExpression(m)
	if !ok {
		errorModel := models.NewError(
			fmt.Sprintf("computed variable '%s' is part of a cyclic dependency", m.selector()),
			pageModel.InputPath,
			m.position,
		)

		documentErrors.AddErrors(&errorModel)
		return
	}

	bootstrap := ""
	for _, p := range index.FindDependentProperties(m) {
		bootstrap = bootstrap + fmt.Sprintf(
			`document.querySelectorAll("[%s]").forEach((__2_element_ref_mod) => __2_element_ref_mod["%s"] = %s);`+"\n",
			p.selector(pageModel), p.propName, p.propAssignment(`"" + (`+expression+`)`),
		)
	}

	if bootstrap == "" {
		return
	}

	pageModel.AppendReactiveRuntime(page.ReactiveRuntimeChunk{
		Variable: m.selector(),
		Content:  bootstrap,
	})
}
