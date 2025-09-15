# 2Web Kit - Virtual DOM

A lightweight virtual DOM.

---

Before you use this package, please first consider if you need a virtual DOM.

A virtual DOM is useful for virtually rendering elements before committing to
the live DOM.

You should only really need this if you are making high-frequency DOM
manipulations that depend on previous DOM state such as `getComputedStyle` or
`getContentBoundingRect` that would typically require a reflow if
read/manipulated in the live DOM.

## Usage

```ts
import { VDomElement, render, vIf } from "@two-web/kit/vdom";

const counterButton = new VDomElement({
  tagName: "button",
  textContent: "0",
  attributes: { "aria-label": "Increment count" },
  events: {
    "click": () => counterButton.textContent++,
  },
  directives: [
    vDomIf(true),
  ],
});

render(document.body, counterButton);
```
