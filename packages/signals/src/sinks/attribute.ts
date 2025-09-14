import type { Signal } from "../signal";

export function attribute<
  Target extends HTMLElement,
  Attribute extends string,
  T,
>(element: Target, name: Attribute, signal: Signal<T>) {
  signal.subscribe((value) => {
    if (value === null || value === undefined) {
      element.removeAttribute(name);
      return;
    }

    element.setAttribute(name, String(value));
  });
}
