# 2Web Tutorial

Beware, 2Web is HIGHLY experimental, and not at all stable.
There are quite a few known bugs.

## Build a "Hello World" application

2Web is based on the "don't pay for what you don't use" philosophy.
If you don't include any reactive code in your website, you will not receive
any JavaScript.

Therefore, a simple "Hello World" website can be written like.

```html
<!DOCTYPE html>
<html>
  <body>
    <h1>Hello World!</h1>
  </body>
</html>
```

If you compile this code in a production build, the only changes that 2web will
make is minifying the supplied code.

_TIP_: You can run a production build with the `--production` flag

We do **not** perform a production build by default, because we believe that
providing readable, code from the development build is important for debugging
your code, understanding what you code does, and integrating with other tools
such as the web browsers development tools.

While we do ship embedded developer tools via the `--dev-tools` flag,
they are not as important as other web frameworks because we have spent a lot
effort ensuring that the web browsers development tools work nicely with the
development build code.
This means that when you get an error, it is a simple browser error that will
match (almost) 1-1 with your source code.

## First reactive app

To be honest, while the "hello world" example was exciting from a technical
standpoint, it doesn't show the exciting reactivity features provided by 2web.

I'll follow the typical reactivity demonstration that all web frameworks do and
demonstrate a simple counter that can increment and a number on the screen.

Because 2Web maintains simplicity as one of its core principles.
I believe that if I can't show the entire website in a single mobile/terminal
window, I have failed in creating a concise framework.
No code snippets will be shown out of context.

```html
<!DOCTYPE html>
<html>
  <body>
    <script compiled>
      $ count = 0;
    </script>

    <h1 *innerText="$count"></h1>
    <button @click="++count">Increment</button>
  </body>
</html>
```

Because 2web is so efficient, I can show you **all** of the compiled code in a
single mobile screen as well.

_NOTE_: This is an dev build. optimized for readability, not performance.

```html
<!DOCTYPE html>
<html>
  <body>
    <h1 data-__2_element_0>0</h1>
    <button onclick="__2_func_0(++__2_var__0)">Increment</button>

    <script>
      let __2_var_0 = 0;
      function __2_func_0(__2_value) {
        document.querySelector("data-__2_element_0")["innerText"] = __2_value;
      }
    </script>
  </body>
</html>
```

The first thing that you will notice is that there are no runtime objects, no
virtual DOM, no zones, no proxies, and no signals.
We are able to achieve reactivity without any runtime overhead.
Changing the variable simply changes the property of the target elements.

An annoying thing about the counter demonstration is that the `$count` compiler
variable still requires a runtime counterpart because its new value is dependent
on its previous value.
If all updates to `$count` did not rely on the previous `$count` value, 2web
would be able to create reactivity without the `__2_var_0` runtime variable.

While it would be possible to remove this runtime variable and boast about a
zero variable runtime, it would greatly hurt performance.

Additionally, all of the first paint content is pre-rendered as SSG content.
You can see in the template above that the initial text of the initial value of
`$count` is inlined as the text `0` in the `data-2-element-0` element.

## "Mustache like" templates

### Why "mustache like"

One of the most common web development tasks is templating text.
I have used a lot of templating systems such as JSX, Svelte htmlx, lit-html, and
Angular's "mustache like" text templating.
While slightly more verbose, I find that mustache templates e.g. `{{ $value }}`
are the easiest to maintain and read.

Templated variables are created using double curly braces as they help
distinguish templated variables against the backdrop of long copy text.

I have found that single curly brackets such as `{ $value }` or `${value}`
make the code harder to read as I (as a programmer) skim over the interpolated
values as grammatical syntax.
Using double curly braces fix this issue because they are visually distinct from
the rest of the text.

As I know, the reasons for chasing "mustache like" text interpolation is a
boring subject, and writing code is more interesting.

### Code example

Expanding upon the "number counter app" we created before, I want to create the
same application using text interpolation becaus having to write out
`*textContent="$value"` is very ugly.

To format text in 2Web, you can simply use "mustache like" brackets similar to
Angular templates.

```html
<script compiled>
  $ count = 0;
</script>

<h1>Current Count: {{ $count }}</h1>
<button @click="++count">Increment</button>
```

The keen eyed of you might have looked at this template and thought to yourself
"this is not a valid html page, is this a snippet of code?"
No.
2Web supports template partials as pages, and can correctly create all of the
boilerplate `DOCTYPE`, `html`, `body`, `head`, `meta`, etc... tags that are
required to create a valid html page.

Therefore, the code above compiles into the following code:

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  </head>

  <body>
    <h1 >Current Count: <span data-__2_element_0>0</span></h1>
    <button onclick="__2_func_0(++__2_var__0)">Increment</button>

    <script>
      let __2_var_0 = 0;
      function __2_func_0(__2_value) {
        document.querySelector("data-__2_element_0")["innerText"] = __2_value;
      }
    </script>
  </body>
</html>
```

As you can see, we've added an additional `<span>` element where the templated
string was located.
This provides more targeted DOM updates, as we don't have to re-render the
parent element or the "Current Count:" text every time the `$count` variable is
updated.

## Document partials

I previously showed that you don't have to write out the HTML boilerplate tags
such as `DOCTYPE`, `html`, `body`, `meta`, etc...

2Web comes with some sane defaults for the html boilerplate.
By default, we set the `DOCTYPE` to `html`, set the `charset` to `UTF-8`, and
set the `viewport` to `width=device-width, initial-scale=1.0`.

Your first thought at seeing this might be the level of control that you give up
when writing document partials.
However, I have created a solution to this specific problem.

A lot of the time, you want to use `<link>` elements to set favicons, set
stylesheets, preload assets, etc... or you might want to set the page title
using a `<title>` tag.
Additionally most major browsers can perform some optimizations if these
elements are in your page's `<head>` tag.

Therefore, any elements that benefit from this sort of optimization, elements
that belong in the `<head>` tag, or are generated by the boilerplate can be
set / overwritten by simply adding them to your document partial.

For example, creating the following document partial.

```html
<title>My Hello World Program</title>
<link rel="stylesheet" href="styles.css" />

<h1>Hello World!</h1>
```

Will expand into the following html code:

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>My Hello World Program</title>
    <link rel="stylesheet" href="styles.css" />
  </head>

  <body>
    <h1>Hello World!</h1>
  </body>
</html>
```

As you can see, the `<title>` and `<link>` elements have been injected into the
`<head>` element.

## Importing components

Because 2Web templates are a superset of HTML, we encounter the same importing
problem as Svelte and Astro where importing components from different files is
a very useful tool, but hard since importing files in HTML doesn't feel quite
right.

E.g. (the following code is an example of what not to do)
`<link rel="component" href="components/footer.component.html" />`
would be a "pure html" way of importing a component from another file.
However, this falls short, and this is non-intuitive for developers coming from
other frameworks.

2Web therefore follows the same import structure that Astro and Svelte make,
where their component imports are "esm like" imports that are evaluated at
compile time.

To import a component in 2Web, you can simply import the component from a
relative path such as:

```html
<script compiled>
  import Footer from "components/footer.component.html";
<script>

<Footer />
```

As you can probably gather from the example above, you can specify the selector
that the component is imported as by changing the default export name.

Additionally, all components do not need an associated closing tag, and can be
terminated by simply adding a `/>` at the end of the element.
