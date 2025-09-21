# "Mustache like" templates

## Why "mustache like"

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

## Code example

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
    <h1>Current Count: <span id="__2_element_0">0</span></h1>
    <button onclick="__2_func_0(++__2_var__0)">Increment</button>

    <script type="module">
      const __2_element_0 = document.getElementById("__2_element_0");
      let __2_var_0 = 0;
      function __2_func_0(__2_value) {
        __2_element_0["innerText"] = __2_value;
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

[**Next**](./4-document-partials.md) (Document partials)
