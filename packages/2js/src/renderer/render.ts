import type { TwoElement } from "../elements/element";
import { change } from "./updates";

export function render(
  target: HTMLElement,
  ...twoElements: TwoElement[]
): void {
  twoElements.forEach((twoEl) => {
    const el = twoEl.toElement();

    change(() => {
      target.appendChild(el);
    });

    twoEl.onUpdate();
  });
}
