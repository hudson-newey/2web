import { TwoElement, iif, classMap } from "../../2js/index.ts";

export const counterButton: any = new TwoElement({
  tagName: "button",
  textContent: "0",
  attributes: {
    id: "counter",
    "aria-label": "Increment count",
  },
  directives: [
    iif(() => parseInt(counterButton.textContent) < 10),
    classMap({
      danger: () => parseInt(counterButton.textContent) >= 5,
    }),
  ],
  events: {
    click: () => counterButton.textContent++,
  },
});
