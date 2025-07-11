# Breaking changes (from HTML)

I have made some minor changes to the HTML format that should not effect you if
you are writing good/proper HTML.

## Code blocks

`<script>`, and `<style>` tags inside of html `<code>` blocks no longer
execute code.
This was done so that you don't have to use `gt;` and `lt;` escape codes inside
of HTML `<code>` blocks.

This allows us to write plain code inside of `<code>` blocks without having to
worry about escape characters.

E.g.

```html
<pre><code>
<h1>My Heading</h1>

<script>
  console.log("Hello World!");
</script>
</code></pre>
```

Running in a normal browser, the inline `<script>` tag would still be executed.
However, I think that this is a horrible developer experience as I just want to
be able to write code inside of `<code>` blocks without having to worry about
html escaping.
Additionally, the subset of developers who are using this is very small, and
I do not expect developers to be writing `<script>` tags that they want to be
executed inside of `<code>` tags.

Therefore, the code above will be compiled to the following HTML.

```html
<pre><code>
&lt;h1&gt;My Heading&lt;/h1&gt;

&lt;script&gt;
  console.log("Hello World!");
&lt;/script&gt;
</pre></code>
```

## Breaking Script Changes

### Script tags have default "module" type

All inline `<script>` tags will automatically support esm imports and will be
lazy-loaded.

### Scripts default to "use strict"

All scripts are run in "strict" mode to encourage good programming practices and
so the browser can make improved performance optimizations.
