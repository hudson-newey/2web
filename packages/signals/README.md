# 2Web Kit - Signals

Lightweight framework agnostic signals that provides **runtime** state tracking.

## Installation

### npm/pnpm/bun/deno node_modules install

```sh
$ npm install @two-web/kit
>
```

### Importing via CDN

If you want to use 2Web signals without using NPM, you can simply add the
following script tag to your documents `<head>`.

```html
<script type="module" src="https://cdn.jsdelivr.net/npm/@two-web/kit/signals"></script>
```

## Usage

```html
<script type="module">
  import { Signal, ComputedSignal, effect } from "@two-web/kit/signals";

  const count = new Signal(0);
  const doubledCount = new ComputedSignal(() => count.value * 2, [count]);

  effect(() => {
    console.log("New value is ${count.value}");
  }, [count]);

  count.subscribe(() => {
    document.getElementById("count-outlet").innerText = value;
  });

  doubledCount.subscribe((value) => {
    document.getElementById("double-count-outlet").innerText = value;
  });
</script>

<button onclick="count.set(count.value + 1)">Increment</button>
<button onclick="count.set(count.value - 1)">Decrement</button>

<output id="count-outlet">0</output>
<output id="double-count-outlet">0</output>
```

### Event Handlers

Note that because the `EventHandler` is a type of signal, it has a `subscribe`
method that will be triggered whenever the event handlers value changes.

```html
<script type="module">
  import { EventHandler } from "@two-web/kit/signals";

  const target = document.getElementById("counter");
  const countHandler = new EventHandler((event, value) => {
    const count = value + 1;

    event.target.textContent = count;
    return count;
  });

  const doubleCount = new ComputedSignal(() => countHandler.value * 2);

  target.addEventListener("click", countHandler);
</script>

<button id="counter">0</button>
```

### Attribute & Property Bindings

```html
<script type="module">
  import { EventHandler, ComputedSignal, attribute } from "@two-web/kit/signals";

  const usernameInput = document.getElementById("username");
  const registerButton = document.getElementById("register-button");

  const username = new EventHandler((event) => event.target.value);
  usernameInput.addEventListener("input", username);

  const isUsernameInvalid = new ComputedSignal(() => {
    username.value.length < 10,
  }, [username]);

  attribute(registerButton, "disabled", isUsernameInvalid);
</script>

<input id="username" />
<button id="register-button" disabled>Register</button>
```

## First Class TypeScript Support

2Web signals have first class structural typing support meaning that if written
correctly, signal state can be known at compile time.

```ts
const greeting = new Signal("Hello");
// ^? Signal<"Hello">;

greeting.set("World!");
// ^? Signal<"World!">;

greeting.update((value) => `Hello ${value}`);
// ^? Signal<"Hello World">
```
