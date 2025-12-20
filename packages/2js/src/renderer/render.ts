import type { TwoElement } from "../elements/element";
import { updateDom } from "../../../_shared/updateDom";

export function render(
  target: HTMLElement,
  ...twoElements: TwoElement[]
): void {
  twoElements.forEach((twoEl) => {
    const el = twoEl.toElement();

    updateDom(() => {
      target.appendChild(el);
    });

    twoEl.onUpdate();
  });
}
