# 2Web Kit - 2JS

A lightweight JavaScript runtime.

---

Before you use this package, please first consider if you need a runtime
library for your website.

A runtime is useful for virtually rendering elements before committing to
the live DOM.

You should only really need this if you are making high-frequency DOM
manipulations that depend on previous DOM state such as `getComputedStyle` or
`getContentBoundingRect` that would typically require a reflow if
read/manipulated in the live DOM.

## Usage

```ts
import { TwoElement, render, iif } from "@two-web/kit/2js";

const counterButton = new TwoElement({
  tagName: "button",
  textContent: "0",
  attributes: { "aria-label": "Increment count" },
  events: {
    "click": () => counterButton.textContent++,
  },
  directives: [
    iif(true),
  ],
});

render(document.body, counterButton);
```
