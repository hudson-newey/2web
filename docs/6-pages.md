# Pages & Routing

2Web uses file-based routing, meaning that the websites routes will be
determined by the structure of your content.

By default, pages are defined as any renderable file that does not have the
`.component` path fragment.

## Layouts

There are two types of layouts, you should pick one and be consistent across
your project.

### Declarative layouts

Declarative layouts must be explicitly consumed, meaning that it's very clear
what layout is being used without looking at the file system path and structure.

This can be thought of like the
[Astro "manual" layouts](https://docs.astro.build/en/basics/layouts/#importing-layouts-manually-mdx)
where to consume a layout, you have a `<slot>` tag, and simply wrap the content
you want the layout to consume.

#### Example

In your `layouts/article.component.html` file, you can create a layout like:

```html
<script compiled>
  import { $props } from "@two-web/compiler";
</script>

<article>
  <h1>{{ $props().title }}</h1>
  <div class="article-content">
    <slot>
      <!-- This content will be rendered if there is no slotted content -->
      Failed to render article
    </slot>
  </div>
</article>

<style>
  .article-content {
    font-family: "Inter", "sans-serif";
    font-size: 1.2rem;
    letter-spacing: 0.2px;
    line-height: 1.6rem;
  }
</style>
```

And this layout can be consumed from within a `blog.html` file like:

```html
<script compiled>
  import BlogLayout from "./layouts/article.component.html";
</script>

<BlogLayout>
  <p>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
  </p>

  <p>
    Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur? Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?
  </p>
</BlogLayout>
```

### Route layouts

Route layouts are _🌠magical🌠_ in that they are special filenames reserved by
the compiler for layouts.

Any children inside the same directory as the `__layout.component.html` file
will use that layout.

Additionally, sub-directories will use the same `__layout.component.html` file
if there is not a conflicting `__layout.component.html` file in the
sub-directory.

This means that `__layout.component.html` files are inherited by sub-pages.

```txt
└─ .
   ├─ __layout.component.html
   └─ article.html
```

Route layouts use the same component format as the declarative layouts, where
content is rendered inside of the `<slot></slot>` elements.

Unlike declarative layouts, consumers of route layouts do **not** need to
explicitly import or template the layout inside of consumers.

[**Next**](./7-control-flow.md) (Template Control Flow)
