# 2 Web

A web framework that compiles straight to lazy loaded HTML, CSS and WASM with
zero runtime overhead.

I am to make all reactivity compiled so that there is zero runtime overhead.

## Resources

- [Tutorial (beta)](./docs/tutorial.md)
- [Theory](./docs/theory.md)

## Basic Counter Example

```html
<script compiled>
$ count = 0;
</script>

<h1 *innerText="$count"></h1>

<button @click="$count = $count + 1">Increment</button>
<button @click="$count = $count - 1">Decrement</button>
```
