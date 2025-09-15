import { VDomElement } from "../../vdom/index.ts";

export const article = new VDomElement({
  tagName: "article",
  children: [
    new VDomElement({
      tagName: "h1",
      textContent: "My Blog",
    }),
    new VDomElement({
      tagName: "article",
      children: [
        new VDomElement({
          tagName: "h2",
          textContent: "Hello World",
        }),
        new VDomElement({
          tagName: "p",
          textContent: "This is my first blog post.",
        }),
      ],
    }),
  ],
});
