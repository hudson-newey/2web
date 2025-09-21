# 2Web Rendering Optimizations

Because 2Web is a compiled web framework, the compiler will typically have
better knowledge about what will change, when it will change, and what will be
effected as a result of the change.

Therefore, in a production build, 2Web can make low-level rendering
optimizations ahead of time.

## CSS Containment (Rasterization Isolation)

Rasterization is the way web browsers paint content to your screen.

Because 2web knows ahead of time what DOM nodes are going to update ahead of
time, the compiler can tell the browser to create a new stacking context for
the updating nodes.

[CSS Containment (mdn)](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_containment/Using_CSS_containment)

Note that because 2Web has a concept of interpolated text range nodes, content
updates can be isolated down to specific words in a sentence.

## `will-change`

While used sparingly, 2Web will attempt to use `will-change` if it deems that an
expensive operation is about to take place.

[will-change (mdn)](https://developer.mozilla.org/en-US/docs/Web/CSS/will-change)

Although `will-change` usage is typically discouraged and warned against in web
development, I have deemed that conservative compiler optimizations are
justifiable because any bugs produced will be consistent and the result of a
logic error, not human error.

Meaning that if implemented correctly, there can only be benefits.

[**Next**](./2-prefetching.md) (Resource Pre-fetching)
