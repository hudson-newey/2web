import { TwoElement, when, classMap, styleMap } from "../../2js/index.ts";

export const counterButton: any = new TwoElement({
  id: "counter",
  className: "big-button",
  tagName: "button",
  textContent: "0",
  attributes: {
    "aria-label": "Increment count",
  },
  directives: [
    when(() => parseInt(counterButton.textContent) < 20),
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
