package models

import (
	"hudson-newey/2web/src/lexer"
)

type ReactiveType int

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

<h1 [innerText $message]></h1>
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
	Static ReactiveType = iota

	// Requires JavaScript to modify the DOM on initial render.
	StaticProperty

	// Requires JavaScript to attach an event listener to the DOM and modify
	// the DOM on event.
	Assignment

	// Requires JavaScript to attach an event listener, keep track of state,
	// and modify the DOM on event.
	Reactive
)

type ReactiveVariable struct {
	Node         *lexer.LexNode[lexer.VarNode]
	Name         string
	InitialValue string
	Bindings     []*ReactiveProperty
}

func (model *ReactiveVariable) AddBinding(property *ReactiveProperty) {
	model.Bindings = append(model.Bindings, property)
}

// TODO: expand this out to assignment and reactive types
// TODO: this should probably cache the type for faster compile times
func (model *ReactiveVariable) Type() ReactiveType {
	if len(model.Bindings) == 0 {
		return Static
	}

	return StaticProperty
}
