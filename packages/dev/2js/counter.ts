import { TwoElement, iif, classMap, styleMap } from "../../2js/index.ts";

export const counterButton: any = new TwoElement({
  tagName: "button",
  textContent: "0",
  attributes: {
    id: "counter",
    "aria-label": "Increment count",
  },
  directives: [
    iif(() => parseInt(counterButton.textContent) < 20),
    classMap({
      danger: () => parseInt(counterButton.textContent) >= 5,
    }),
    styleMap({
      fontSize: () =>
        parseInt(counterButton.textContent) >= 10 ? "6rem" : "32px",
    }),
  ],
  events: {
    click: () => counterButton.textContent++,
  },
});
