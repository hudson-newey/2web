import { TwoElement } from "../../2js/index.ts";

export const article = new TwoElement({
  tagName: "article",
  children: [
    new TwoElement({
      tagName: "h1",
      textContent: "My Blog",
    }),
    new TwoElement({
      tagName: "article",
      children: [
        new TwoElement({
          tagName: "h2",
          textContent: "Hello World",
        }),
        new TwoElement({
          tagName: "p",
          textContent: "This is my first blog post.",
        }),
      ],
    }),
  ],
});
