# Component Authoring & Usage

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
