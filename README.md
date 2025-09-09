# 2 Web

A web framework that compiles straight to lazy loaded HTML, CSS and JS with
near-zero runtime overhead.

I am to make all reactivity compiled so that there is minimal runtime overhead.

## Resources

- [Tutorial (beta)](./docs/README.md)
- [CLI Documentation](./cli/README.md)

## Basic Counter Example

```html
<script>
  $ count = 0;
</script>

<h1>{{ $count }}</h1>

<button @click="$count = $count + 1">Increment</button>
<button @click="$count = $count - 1">Decrement</button>
```

## Quick Start

Install the 2web cli.

```sh
$ npm install -g @two-web/cli
>
```

Generate a new 2web project.

```sh
$ 2web new <project_name>
>
```
