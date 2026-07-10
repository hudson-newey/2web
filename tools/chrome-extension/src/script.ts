import { render, TwoElement } from "@two-web/kit/2js";

const app = new TwoElement({
  tagName: "main",
  textContent: "Hello World!",
});

render(document.getElementById("output-element")!, app);
