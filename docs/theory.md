# Theory

With 2Web, you get what you write, if you write the following code

```html
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    Hello
  </body>
</html>
```

We will not inject any additional reactivity into the page

Your first instict might be "well what if we want the word 'Hello' to be
reactive?".

```html
<script compiled>
  $ text = "Hello"
</script>

<template>
  <span>{{ $text }}</span>
  <button @click="text = 'world'">Change</button>
</template>
```

2Web will compile the above code to the following (unminified) code.

```html
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <span id="out-element">Hello</span>
    <button onclick="change0()"></button>
  </body>

  <script>
    function change0() {
      document.getElementById("out-element").innerText = "world";
    }
  </script>
</html>
```

Want to create an incrementing number?

```html
<script compiled>
  $ number = 0
</script>

<template>
  <span>{{ $number }}</span>
  <button @click="number++">Increment</button>
</template>
```

2Web will compile the above code to the following (unminified) code.

```html
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <span id="out-element">0</span>
    <button onclick="increment0()"></button>
  </body>

  <script>
    let $number = 0;

    function increment0() {
      $number++;
      document.getElementById("out-element").innerText = $number;
    }
  </script>
</html>
```

As you can see, the compiled code does include a single JS variable to store
the reactive data, but there is no runtime overhead.

## Computed reactivity

Web libraries are starting to move away from VDOM's and towards property
tracking reactivity.
However, these libraries use runtime signals to track dependencies.

A simple version of fine grained property tracking reactivity can be
implemented with the following runtime code.

```ts
type Subscriber = () => void;

// this class would typically use a WeakMap to store the subscribers, it would
// typically have a way to unsubscribe, and would also use Proxy objects for
// property subscriptions.
class Signal<T> {
    #value: T;
    #subscribers: Subscriber[] = [];

    public subscribe(callback: Subscriber) {
        this.#subscribers.push(callback);
    }

    public set(value: <T>) {
        this.#value = value;

        for (const subscriber of this.#subscribers) {
            subscriber(value);
        }
    }

    public get() {
        return this.#value;
    }
}

function reactiveProperty<T extends HTMLELement>(
    target: T,
    key: keyof T,
    signal: Signal,
) {
    signal.subscribe((value) => {
        target[key] = value;
    });
}

const signal = new Signal(0);

signal.set(1);

document.setInnerHTML = `
    <button onclick="signal.set(signal.get() + 1)">+ 1</button>
    <span onload="reactiveProperty(event.target, 'innerText', signal)"></span>
`;
```

A lot of the frameworks you love such as Svelte, Solid, and Vue vapor mode
use compilers so that you can have shorthands for the reactiveProperty
function.
However, this is not true compiler reactivity.
It is just compiler shorthands for a very fast runtime reactivity system.

Using a compiler in these situations can be helpful because you can tree shake
and optimize some of the code.
But compiler optimizations from these frameworks typically do not optimize
away signal based reactivity.

However, this is sub optimal because it still has a Signal class, and it still
has a reactiveProperty function.

## Compiler based computed reactivity

As mentioned previously, these frameworks typically use signals so that
computed properties can be reactive.

However, it is possible (although hard) to create computed reactivity at
compile time, removing the need for signals.

e.g. The following code is a simple example with templated computed values

````html
```html
<script compiled>
  $ number = 0
</script>

<template>
  <span>{{ $number * 2 }}</span>
  <button @click="number++">Increment</button>
</template>
````

This is easy for the compiler, because it can just inline the number
calculation into the onclick function.

```html
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <span id="out-element">0</span>
    <button onclick="increment0()"></button>
  </body>

  <script>
    let $number = 0;

    function increment0() {
      number++;
      document.getElementById("out-element").innerText = $number * 2;
    }
  </script>
</html>
```

However, more complex examples are harder.

```html
<script compiled>
  $ itemId = 0

  // fetch the number from the server
  $ itemTitle = (async () => {
      const response = await fetch(`https://example.com/${$itemId}`);
      const responseBody = response.json();
      return responseBody.title;
  })();
</script>

<template>
  <span>{{ $itemTitle }}</span>
  <button @click="itemId++">Change Item Id</button>
</template>
```

As you can probably notice, we can first inline the callback

```html
<script compile>
  $ itemId = 0
</script>

<template>
  <span
    >{{ (() => { const response = await fetch(`https://example.com/${$itemId}`);
    const responseBody = response.json(); return responseBody.title; }) }}</span
  >
  <button @click="itemId++"></button>
</template>
```

Then we can simply inline the function to the compiled code

```html
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <span id="out-element"></span>
    <button onclick="increment0()"></button>
  </body>

  <script>
    let $itemId = 0;

    function increment0() {
      itemId++;
      document.getElementById("out-element").innerText = (async () => {
        const response = await fetch(`https://example.com/${$itemId}`);
        const responseBody = response.json();
        return responseBody.title;
      })();
    }
  </script>
</html>
```

And this is exactly what we do

See. It's not as hard as you think.

What about objects and reactivity? Wouldn't I need to use a Proxy object? No.

Since we use a compiler, we know all of the properties that will be accessed
at runtime.
I can easily destructure the object and make each individual property reactive.
