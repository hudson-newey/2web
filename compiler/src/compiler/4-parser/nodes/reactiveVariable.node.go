package nodes

import (
	"fmt"
	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/2-lexer/lexeme"
	"hudson-newey/2web/src/compiler/4-parser/scanners"
	"hudson-newey/2web/src/content/css"
	"hudson-newey/2web/src/content/javascript"
	"hudson-newey/2web/src/content/page"
	twoscript "hudson-newey/2web/src/content/twoScript"
	"slices"
	"strings"

	"github.com/hudson-newey/2web/_shared/lists"
	"github.com/hudson-newey/2web/_shared/logger"
)

func NewreactiveVariableNode(lexNodes []*lexer.V2LexNode) *reactiveVariableNode {
	variableName, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 1)
	if err != nil {
		panic(err)
	}

	initialValue, err := scanners.NthToken(lexNodes, lexeme.CompiledScriptSource, 2)
	if err != nil {
		panic(err)
	}

	return &reactiveVariableNode{
		variableName: variableName.Content,
		initialValue: initialValue.Content,
	}
}

type reactiveVariableNode struct {
	variableName string
	initialValue string
	children     AbstractSyntaxTree
}

func (m *reactiveVariableNode) Type() string {
	return "reactiveVariableNode"
}

func (m *reactiveVariableNode) Children() AbstractSyntaxTree {
	return m.children
}

func (m *reactiveVariableNode) Content(page *page.Page, _ast AbstractSyntaxTree) NodeContent {
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

func (m *reactiveVariableNode) dependentProps(ast AbstractSyntaxTree) []*reactivePropertyNode {
	return lists.Filter(ast.reactiveProperties(), func(x *reactivePropertyNode) bool {
		return slices.Contains(x.reactiveVariableDeps(ast), m)
	})
}

func (m *reactiveVariableNode) dependentEvents(ast AbstractSyntaxTree) []*reactiveEventNode {
	return lists.Filter(ast.reactiveEvents(), func(x *reactiveEventNode) bool {
		return slices.Contains(x.reactiveVariableDeps(ast), m) || x.reactiveVariableSink(ast) == m
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

<h1 *innerText="$message"></h1>
```

As you can see from the example, the variable is not really reactive, but it
does require a <script> tag to modify the DOM on initial render.

The compiled code will look something like this:

```html
<h1 id="_0">Hello!</h1>
<script>

	document.addEventListener("DOMContentLoaded", () => {
	    document.getElementById("_0").innerText = "Hello!";
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
	// does not require any JavaScript. Can be inlined at compile time.
	static reactivityLevel = iota

	// Requires JavaScript to modify the DOM on initial render.
	staticProperty

	// Requires JavaScript to attach an event listener to the DOM and modify
	// the DOM on event.
	assignment

	// Requires JavaScript to attach an event listener, keep track of state,
	// and modify the DOM on event.
	reactive
)

func (m *reactiveVariableNode) reactivityLevel(ast AbstractSyntaxTree) reactivityLevel {
	events := m.dependentEvents(ast)
	for i := range events {
		event := events[i]
		// If the assignment expression uses the same variable that it's
		// assigning to, we need to have a runtime variable to track state.
		// e.g. think of a counting number
		if strings.Contains(event.assignmentExpr, m.selector()) {
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
	props := m.dependentProps(ast)
	if len(props) > 0 {
		return staticProperty
	}

	return static
}
