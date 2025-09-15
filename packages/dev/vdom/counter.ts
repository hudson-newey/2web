import { VDomElement, vDomIf } from "../../vdom/index.ts";

export const counterButton: any = new VDomElement({
  tagName: "button",
  textContent: "0",
  attributes: {
    id: "counter",
    "aria-label": "Increment count",
  },
  directives: [vDomIf(() => parseInt(counterButton.textContent) < 10)],
  events: {
    "click": () => counterButton.textContent++,
  },
});

