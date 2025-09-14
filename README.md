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

## Structure

### [Compiler](./compiler/)

The 2Web compiler can optimize existing HTML, JavaScript/TypeScript, and
CSS/SCSS/SASS code while also providing a "fixed" syntax of HTML.

#### Optimizations

1. The compiler can recognize that a resource is not on the hot-path, and will
   lazy load the resource.
2. The compiler can identify dynamic parts of the page and automatically apply
   `contains` and `will-change` optimizations.
3. (**TODO**) Automatically create pages for dynamic routes using the
   `[param_name]` file name syntax.

#### Improvements

1. First-class TypeScript support in `<script>` tags
2. (**TODO**) First-class SASS/SCSS support in `<style>` and `<link>` tags
3. (**TODO**) All elements can now self-close (e.g. `<script src="my-script.ts" />`)
4. (**TODO**) You can escape HTML tokens (e.g. `<`) by using back slashes
   (e.g. `/<` will display a `<` character).
5. Tags inside of `<code>` blocks no longer need to use HTML escape codes.
6. You can import `.html` files as components.
7. New JavaScript reactive variable declaration syntax using `$` keyword.
8. Create web pages without HTML boilerplate
9. Create web pages using Markdown

### [2web/kit](./packages/)

The 2Web kit is a package containing framework agnostic helpers to build fast
and robust web applications.

### [CLI](./cli/)

Provides a standardized (opinionated) way to set up a project.

The CLI also contains a lot of templates and generators that can be used to
easily add functionality such as SSR, database, load balancers, etc... to your
application with one command.

### VsCode Extension

A vscode extension to integrate with the 2web cli.
